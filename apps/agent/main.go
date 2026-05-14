package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"

	rcmasvc "agent/internal/service"
)

// usage prints a short help message to stderr.
func usage() {
	fmt.Fprintln(os.Stderr, `RCMA Agent Usage:
  rcma                  Run in console/debug mode (no service required)
  rcma install          Install the Windows service
  rcma uninstall        Uninstall the Windows service
  rcma start            Start the installed service
  rcma stop             Stop the installed service
  rcma restart          Restart the installed service
  rcma status           Print current service status
`)
}

func main() {
	prg, err := rcmasvc.New()
	if err != nil {
		log.Fatalf("Failed to initialise agent: %v", err)
	}

	svcConfig := rcmasvc.ServiceConfig()
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	// ── CLI command handling ──────────────────────────────────────────────────
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install", "uninstall", "start", "stop", "restart":
			if err := service.Control(svc, cmd); err != nil {
				log.Fatalf("service control %q failed: %v", cmd, err)
			}
			fmt.Printf("Service %q executed successfully.\n", cmd)
			return

		case "status":
			status, err := svc.Status()
			if err != nil {
				log.Fatalf("Cannot query service status: %v", err)
			}
			switch status {
			case service.StatusRunning:
				fmt.Println("Status: Running")
			case service.StatusStopped:
				fmt.Println("Status: Stopped")
			default:
				fmt.Println("Status: Unknown")
			}
			return

		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %q\n\n", cmd)
			usage()
			os.Exit(1)
		}
	}

	// ── Run ───────────────────────────────────────────────────────────────────
	// If launched by the SCM the interactive check is false and service.Run
	// hands control to the SCM.  If launched from a terminal (debug mode) it
	// calls Start/Stop directly so the agent runs in the console, making it
	// easy to attach a debugger or watch log output in real time.
	if err := svc.Run(); err != nil {
		log.Fatalf("Service run failed: %v", err)
	}
}
