package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"guiinstaller/internal/constants"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/svc/mgr"
)

type Config struct {
	LANIP     string `json:"lan_ip"`
	Port      int    `json:"port"`
	AuthToken string `json:"auth_token"`
}

// DownloadAgent downloads a file from URL to destination with progress callback
// progressCallback is called with the progress percentage (0-100)
func DownloadAgent(url, dest string, progressCallback func(float64)) error {
	// 1. Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ERROR E20 - failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ERROR E21 - bad status: %s", resp.Status)
	}

	// ✅ Validate PE header (MZ) without draining the body
	header := make([]byte, 2)
	_, err = io.ReadFull(resp.Body, header)
	if err != nil {
		return fmt.Errorf("ERROR E21 - failed to read file header: %v", err)
	}
	if header[0] != 'M' || header[1] != 'Z' {
		return fmt.Errorf("ERROR E21 - downloaded file is not a valid Windows executable (got: %x %x)", header[0], header[1])
	}

	// Stitch the 2 bytes back so we don't lose them during copy
	fullBody := io.MultiReader(bytes.NewReader(header), resp.Body)

	// 2. Create directory
	if err := os.MkdirAll(constants.RCMA_BINARY_DIR_PATH, 0755); err != nil {
		return fmt.Errorf("ERROR E22 - Failed to create agent binary directory 'AppData/Local/RCMA': %v", err)
	}

	// 3. Create destination file
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("ERROR E22 - failed to create file: %v", err)
	}
	defer out.Close()

	// 4. Get total file size
	totalSize := resp.ContentLength

	// 5. Copy with progress tracking
	buffer := make([]byte, 32*1024)
	var downloadedSize int64 = 0

	for {
		n, err := fullBody.Read(buffer)
		if n > 0 {
			_, writeErr := out.Write(buffer[:n])
			if writeErr != nil {
				return fmt.Errorf("ERROR E23 - failed to save file: %v", writeErr)
			}
			downloadedSize += int64(n)

			if totalSize > 0 && progressCallback != nil {
				progress := float64(downloadedSize) / float64(totalSize) * 100
				progressCallback(progress)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ERROR E24 - failed to download: %v", err)
		}
	}

	if progressCallback != nil {
		progressCallback(100)
	}

	return nil
}

func SetupForAgent() error {

	// 2) Create dir C:\ProgramData\RCMA and config.json file at there
	if err := os.MkdirAll(constants.CONFIG_DIR, 0755); err != nil {
		return fmt.Errorf("ERROR E26 - Failed to create config directory: %v", err)
	}

	f, err := os.Create(constants.CONFIG_FILE_PATH)
	if err != nil {
		return fmt.Errorf("ERROR E27 - Failed to create config file: %v", err)
	}

	defer f.Close()

	// 3) Gather important resources to store in config.json
	addr, lanError := GetLANIP()
	if lanError != nil {
		return fmt.Errorf("ERROR E28 - Error on getting LANIP: %v", lanError)
	}

	port, portError := PickPort()
	if portError != nil {
		return fmt.Errorf("ERROR E29 - Error on getting PORT: %v", portError)
	}

	token, tokenError := GenerateSecureToken(32)
	if tokenError != nil {
		return fmt.Errorf("ERROR E30 - Error on generating token: %v", tokenError)
	}

	// 4) Write down resources in config.json
	cfg := Config{LANIP: addr, Port: port, AuthToken: token}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("ERROR E31 - Failed to write config: %v", err)
	}

	return nil
}

func DeleteDir(dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return fmt.Errorf("ERROR E32 - Directory does not exist: %s\n", dest)
	}

	abPath, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("ERROR E33 - Error resolving absolute path %v\n", err)
	}
	if abPath == "/" || abPath == "C:\\" {
		return fmt.Errorf("ERROR E34 - Refusing to delete root directory for safety.")
	}
	err = os.RemoveAll(dest)
	if err != nil {
		return fmt.Errorf("ERROR E35 - Error delete directory: %v\n", err)
	}

	return nil
}

func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", fmt.Errorf("ERROR E36 - byte length must be positive")
	}

	// Create a byte slice to hold random data
	tokenBytes := make([]byte, byteLength)

	// Fill the slice with secure random bytes
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("ERROR E37 - failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(tokenBytes), nil
}

func GetLANIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("ERROR E38 - failed to dial network: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func PickPort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("ERROR E39 - failed to listen on port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port, nil
}

func RegisterService(agentPath string, agentName string, displayName string) error {

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("ERROR E40 - Error on registering agent as service: %v", err)
	}
	defer m.Disconnect()

	s, err := m.CreateService(agentName, agentPath, mgr.Config{
		DisplayName: displayName,
		StartType:   mgr.StartAutomatic,
	})

	if err != nil {
		return fmt.Errorf("ERROR E41 - Error on creating service: %v", err)
	}

	defer s.Close()

	fmt.Println("Service registered!")
	return nil
}

func IsServiceExist(agentName string) (bool, error) {

	m, err := mgr.Connect()
	if err != nil {
		return false, fmt.Errorf("ERROR E42 - failed to connect to service manager: %v", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(agentName)
	if err != nil {
		if errors.Is(err, syscall.Errno(1060)) {
			return false, nil
		} else {
			return false, fmt.Errorf("ERROR E44 - Error on opening service: %v\n", err)
		}
	} else if s == nil {
		return false, nil
	}
	defer s.Close()

	return true, nil
}
