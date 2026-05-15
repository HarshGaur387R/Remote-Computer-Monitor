package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kardianos/service"

	"agent/internal/agent"
	"agent/internal/logger"
)

const (
	rcmaDirPath   = `C:\ProgramData\RCMA`
	configPath    = rcmaDirPath + `\config.json`
	agentLogsPath = rcmaDirPath + `\agent_logs.ndjson`
)

// Program implements service.Interface for kardianos/service.
type Program struct {
	log    *logger.Logger
	server *http.Server
}

// New creates a Program. The logger is initialised here so it is available
// both in service mode and in console/debug mode.
func New() (*Program, error) {
	l, err := logger.New(rcmaDirPath)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	return &Program{log: l}, nil
}

// Start is called by the service manager to start the program.
// It must not block; long-running work goes in a goroutine.
func (p *Program) Start(s service.Service) error {
	p.log.Info("Service Start called — launching agent goroutine")
	go p.run()
	return nil
}

// Stop is called by the service manager to stop the program gracefully.
func (p *Program) Stop(s service.Service) error {
	p.log.Info("Service Stop called — initiating graceful shutdown")
	if p.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.server.Shutdown(ctx); err != nil {
			p.log.Errorf("HTTP server shutdown error: %v", err)
			return err
		}
	}
	p.log.Info("Service stopped cleanly")
	p.log.Close()
	return nil
}

// run contains the actual agent startup logic.
func (p *Program) run() {
	p.log.Info("Agent starting")

	cfg, err := agent.LoadConfig(configPath)
	if err != nil {
		p.log.Errorf("Failed to load config: %v", err)
		return
	}
	p.log.Infof("Config loaded — binding on %s:%d", cfg.LANIP, cfg.Port)

	router := agent.BuildRouter(cfg.AuthToken)

	addr := fmt.Sprintf("%s:%d", cfg.LANIP, cfg.Port)
	p.server = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	p.log.Infof("HTTP server listening on http://%s", addr)
	if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		p.log.Errorf("HTTP server fatal error: %v", err)
	}
}

// ServiceConfig returns the kardianos/service configuration for RCMA.
func ServiceConfig() *service.Config {
	return &service.Config{
		Name:        "rcma",
		DisplayName: "RCMA Remote Monitoring Agent",
		Description: "Exposes system metrics over HTTP for RCMA.",
	}
}
