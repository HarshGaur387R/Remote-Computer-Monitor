package main

import (
	rcmasvc "agent/internal/service"
	"github.com/kardianos/service"
	"log"
)

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

	// ── Run ───────────────────────────────────────────────────────────────────
	// If launched by the SCM the interactive check is false and service.Run
	// hands control to the SCM.  If launched from a terminal (debug mode) it
	// calls Start/Stop directly so the agent runs in the console, making it
	// easy to attach a debugger or watch log output in real time.
	if err := svc.Run(); err != nil {
		log.Fatalf("Service run failed: %v", err)
	}
}