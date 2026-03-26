package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Config struct {
	LANIP     string `json:"lan_ip"`
	Port      int    `json:"port"`
	AuthToken string `json:"auth_token"`
}

func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", fmt.Errorf("byte length must be positive")
	}

	// Create a byte slice to hold random data
	tokenBytes := make([]byte, byteLength)

	// Fill the slice with secure random bytes
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(tokenBytes), nil
}

func getLANIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func pickPort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port, nil
}

func main() {

	configDir := `C:\ProgramData\RCMAgent`
	configFile := configDir + `\config.json`

	addr, lanError := getLANIP()
	if lanError != nil {
		log.Fatalf("Error on getting LANIP: %v", lanError)
	}

	port, portError := pickPort()
	if portError != nil {
		log.Fatalf("Error on getting PORT: %v", portError)
	}

	token, tokenError := GenerateSecureToken(32)
	if tokenError != nil {
		log.Fatalf("Error on generating token: %v", tokenError)
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	f, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}

	defer f.Close()

	cfg := Config{LANIP: addr, Port: port, AuthToken: token}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		log.Fatalf("Failed to write config: %v", err)
	}

	fmt.Println("Hello from installer")
}
