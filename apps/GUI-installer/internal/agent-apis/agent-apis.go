package agentapis

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func Start_agent(agentName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("ERROR - failed to connect to service manager: %v", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(agentName)
	if err != nil {
		if errors.Is(err, syscall.Errno(1060)) {
			return fmt.Errorf("ERROR E100 - Service does not exist")
		} else {
			return fmt.Errorf("ERROR E101 - Error on opening service: %v\n", err)
		}
	}
	defer s.Close()

	startErr := s.Start()
	if startErr != nil {
		return fmt.Errorf("ERROR E102 - Error on starting agent: %v\n", startErr)
	}

	return nil
}

func Stop_agent(agentName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("ERROR - failed to connect to service manager: %v", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(agentName)
	if err != nil {
		if errors.Is(err, syscall.Errno(1060)) {
			return fmt.Errorf("ERROR E103 - Service does not exist")
		} else {
			return fmt.Errorf("ERROR E104 - Error on opening service: %v", err)
		}
	}
	defer s.Close()

	_, err = s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("ERROR E105 - Error on stopping agent: %v", err)
	}

	return nil
}

func Restart_agent(agentName string) error {
	if err := Stop_agent(agentName); err != nil {
		return fmt.Errorf("ERROR E106 - Error on stopping agent during restart: %v", err)
	}

	// Add a small delay to ensure service is fully stopped
	time.Sleep(500 * time.Millisecond)

	if err := Start_agent(agentName); err != nil {
		return fmt.Errorf("ERROR E107 - Error on starting agent during restart: %v", err)
	}

	return nil
}

func IsRunning(agentName string) (bool, error) {
	m, err := mgr.Connect()
	if err != nil {
		return false, fmt.Errorf("ERROR - failed to connect to service manager: %v", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(agentName)
	if err != nil {
		if errors.Is(err, syscall.Errno(1060)) {
			return false, fmt.Errorf("ERROR E108 - Service does not exist")
		} else {
			return false, fmt.Errorf("ERROR E109 - Error on opening service: %v", err)
		}
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return false, fmt.Errorf("ERROR E110 - Error querying service status: %v", err)
	}

	return status.State == svc.Running, nil
}
