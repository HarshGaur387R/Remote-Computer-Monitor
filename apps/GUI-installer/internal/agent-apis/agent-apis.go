package agentapis

import (
	"guiinstaller/internal/constants"
	"guiinstaller/internal/utils"

	"github.com/kardianos/service"
)

func serviceHandle(agentBinaryPath string) (service.Service, error) {
	// The GUI creates the same service.Config the agent uses.
	// kardianos/service uses this only to talk to SCM — it does NOT run the binary.
	cfg := &service.Config{
		Name:        constants.AGENT_NAME,
		DisplayName: "RCMA Remote Monitoring Agent",
		Description: "Exposes system metrics over HTTP for RCMA.",
		Executable:  agentBinaryPath,
	}
	// nil program — GUI never calls Start/Stop on the Program interface,
	// only on the SCM handle
	return service.New(nil, cfg)
}

func IsServiceExist(name string) (bool, error) {
	svc, err := serviceHandle(constants.RCMA_BINARY_PATH)
	if err != nil {
		return false, err
	}
	_, err = svc.Status()
	if err == service.ErrNotInstalled {
		return false, nil
	}
	return err == nil, err
}

func RegisterService(binaryPath, _, _ string) error {
	if err := utils.SetupForAgent(); err != nil {
		return err
	}
	svc, err := serviceHandle(binaryPath)
	if err != nil {
		return err
	}
	return svc.Install()
}

func UnregisterService(binaryPath string) error {
	svc, err := serviceHandle(binaryPath)
	if err != nil {
		return err
	}
	return svc.Uninstall()
}

// internal/agent-apis/apis.go

func IsRunning(name string) (bool, error) {
	svc, err := serviceHandle(constants.RCMA_BINARY_PATH)
	if err != nil {
		return false, err
	}
	status, err := svc.Status()
	return status == service.StatusRunning, err
}

func Start_agent(name string) error {
	svc, err := serviceHandle(constants.RCMA_BINARY_PATH)
	if err != nil {
		return err
	}
	return service.Control(svc, "start")
}

func Stop_agent(name string) error {
	svc, err := serviceHandle(constants.RCMA_BINARY_PATH)
	if err != nil {
		return err
	}
	return service.Control(svc, "stop")
}
