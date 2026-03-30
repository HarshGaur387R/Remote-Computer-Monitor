package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
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

func printHelp() {
	fmt.Println()
	fmt.Println("RCMAgent Installer - Command Reference")
	fmt.Println("======================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  rcmaI <command>")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  start    Launches the agent service. If already running, no action is taken. If not running,")
	fmt.Println("           the installer checks whether the agent is installed. If missing, the latest release")
	fmt.Println("           is downloaded and installed before starting the service.")
	fmt.Println()
	fmt.Println("  stop     Stops the agent service if it is currently running.")
	fmt.Println()
	fmt.Println("  restart  Restarts the agent service by stopping it first and then starting it again.")
	fmt.Println()
	fmt.Println("  status   Reports whether the agent service is running. Also indicates if the installed")
	fmt.Println("           agent version is outdated compared to the latest release.")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/harshgaur/rcm_monorepo")
	fmt.Println()
}

// Returns true if the service is running, false otherwise
func isServiceRunning(serviceName string) (bool, error) {
	out, err := exec.Command("sc", "query", serviceName).Output()
	if err != nil {
		return false, err
	}
	result := string(out)
	// Look for "RUNNING" in the output
	return strings.Contains(result, "RUNNING"), nil
}

// Returns true if the service exists, false otherwise
func serviceExists(serviceName string) (bool, error) {
	out, err := exec.Command("sc", "query", serviceName).CombinedOutput()
	if err != nil {
		// If the error output contains "does not exist", we know it’s missing
		if strings.Contains(string(out), "does not exist") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func startAgent() {
	configDir := `C:\ProgramData\RCMAgent`
	configFile := configDir + `\config.json`
	serviceName := "RCMAgentService"

	// Chekcking if service exist.
	isServiceExists, err := serviceExists(serviceName)
	if err != nil {
		fmt.Println("Error on finding service:", err)
		return
	}
	if !isServiceExists {
		fmt.Println("Service not found, Try to install it again. run `rcmai --update`")
		return
	}

	// Checking if service is running
	running, err := isServiceRunning(serviceName)
	if err != nil {
		fmt.Println("Error checking service:", err)
		return
	}
	if running {
		fmt.Println("RCMAgentService is already running")
		return
	} else {
		fmt.Println("RCMAgentService is not running")
	}

	return

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

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No command is given. Run program with --help\nExample: rcmai --help")
		return
	}

	switch os.Args[1] {
	case "--help":
		printHelp()
	case "start": // start service + show QR
		startAgent()
	case "stop": // stop service
		fmt.Println("Stop is pressed")
	case "status": // query service
		fmt.Println("Status is pressed")
	case "restart": // stop + start
		fmt.Println("Restart is pressed")
	}
}
