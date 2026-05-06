package main

import (
	"fmt"
	agentapis "guiinstaller/internal/agent-apis"
	"guiinstaller/internal/constants"
	"guiinstaller/internal/utils"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func startAgent(statusState binding.Bool, parent fyne.Window) {
	{

		errOnSetup := utils.SetupForAgent()
		agentBinaryInfo, err := os.Stat(constants.RCMA_BINARY_PATH)

		if err != nil {
			dialog.ShowError(fmt.Errorf("ES1 Error on starting agent: %v", err), parent)
			return
		}

		if agentBinaryInfo == nil {
			dialog.ShowError(fmt.Errorf("ES2 Error on starting agent: %v", err), parent)
			return
		}

		if errOnSetup != nil {
			dialog.ShowError(errOnSetup, parent)
			return
		}

		exist, err := utils.IsServiceExist(constants.AGENT_NAME)
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if !exist {
			err := utils.RegisterService(constants.RCMA_BINARY_PATH, constants.AGENT_NAME, constants.AGENT_DISPLAY_NAME)
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
		}

		running, err := agentapis.IsRunning(constants.AGENT_NAME)

		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if running {
			dialog.ShowInformation("Agent is already running", "Click stop if you want to stop agent from running", parent)
			return
		} else {
			err := agentapis.Start_agent(constants.AGENT_NAME)

			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
		}

		statusState.Set(true)
	}
}

func stopAgent(statusState binding.Bool, parent fyne.Window) {
	// Pause button action

	exist, err := utils.IsServiceExist(constants.AGENT_NAME)
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	if exist {
		running, err := agentapis.IsRunning(constants.AGENT_NAME)

		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if running {
			err := agentapis.Stop_agent(constants.AGENT_NAME)
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			statusState.Set(false)
			return
		} else {
			dialog.ShowInformation("Agent is not running", "Agent is already in stop state, to run it press start.", parent)
			return
		}
	} else {
		dialog.ShowError(fmt.Errorf("Agent does not exist"), parent)
		statusState.Set(false)
		return
	}
}

func restartAgent(statusState binding.Bool, parent fyne.Window) {
	// Restart button action
	exist, err := utils.IsServiceExist(constants.AGENT_NAME)
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	if exist {
		isRunning, err := agentapis.IsRunning(constants.AGENT_NAME)

		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if isRunning {
			stopAgent(statusState, parent)

			// Small delay to ensure service is fully stopped
			time.Sleep(500 * time.Millisecond)

			startAgent(statusState, parent)
			statusState.Set(true)
			return

		} else {
			startAgent(statusState, parent)
			statusState.Set(true)
			return
		}

	}
}

func AgentTab(parent fyne.Window) fyne.CanvasObject {

	// ==================== STATES ====================

	statusState := binding.NewBool()
	statusState.Set(false)

	// ==================== HEADER ====================
	header := canvas.NewText("Agent Control", color.RGBA{
		R: 0,
		G: 0,
		B: 200,
		A: 255,
	})
	header.TextStyle.Bold = true
	header.TextSize = 32

	headerContainer := container.NewCenter(header)

	// Create a container to hold the dynamic content
	contentBox := container.NewVBox()

	UpdateUI := func() {

		service_exist, err := utils.IsServiceExist(constants.AGENT_NAME)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error on checking if agent exist: %v", err), parent)
			statusState.Set(false)
		}

		if service_exist {
			is_running, err := agentapis.IsRunning(constants.AGENT_NAME)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Error on checking if agent is running: %v", err), parent)
				statusState.Set(false)
			}
			statusState.Set(is_running)
		} else {
			statusState.Set(false)
		}

		// ==================== STATES =========================
		statusValue, _ := statusState.Get()
		contentBox.Objects = nil

		// ==================== CONTROL BAR ====================
		playBtn := widget.NewButton("▶ Start", func() { startAgent(statusState, parent) })
		pauseBtn := widget.NewButton("⏸ Stop", func() { stopAgent(statusState, parent) })
		restartBtn := widget.NewButton("🔄 Restart", func() { restartAgent(statusState, parent) })

		// Determine status text and color
		statusText := "UnActive"
		statusColor := color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255, // Red for inactive
		}
		if statusValue {
			statusText = "Active"
			statusColor = color.RGBA{
				R: 0,
				G: 255,
				B: 0,
				A: 255, // Green for active
			}
		}

		// Create status display
		statusTextCanvas := canvas.NewText(statusText, statusColor)
		statusTextCanvas.TextSize = 14
		statusContainer := container.NewHBox(
			widget.NewLabel("Status:"),
			statusTextCanvas,
		)

		// Control buttons container
		controlButtonsContainer := container.NewHBox(
			playBtn,
			pauseBtn,
			restartBtn,
		)

		// Combine control bar
		controlBar := container.NewVBox(
			statusContainer,
			controlButtonsContainer,
		)

		controlBarWithBorder := container.NewBorder(
			nil, nil, nil, nil,
			controlBar,
		)

		contentBox.Add(controlBarWithBorder)
		contentBox.Refresh()
	}
	// ==================== LOGS DISPLAY ====================
	logTitle := widget.NewLabel("Agent Logs")
	logTitle.TextStyle.Bold = true

	// Create green text for logs with proper line breaks
	greenColor := color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255, // Green text
	}

	logsInner := container.NewVBox(
		canvas.NewText("Logs will appear here...", greenColor),
		canvas.NewText("", greenColor),
		canvas.NewText("- Agent initialized", greenColor),
		canvas.NewText("- Ready for connections", greenColor),
	)

	// Black background rectangle
	blackBg := canvas.NewRectangle(color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255, // Black background
	})

	// Scrollable container for logs
	logScroll := container.NewScroll(logsInner)
	logScroll.SetMinSize(fyne.NewSize(600, 400))

	// Combine background and logs
	logsContent := container.NewStack(
		blackBg,
		container.NewVBox(
			logScroll,
		),
	)

	logsContainer := container.NewVBox(
		logTitle,
		logsContent,
	)

	statusState.AddListener(binding.NewDataListener(UpdateUI))

	UpdateUI()

	// ==================== MAIN LAYOUT ====================
	screen := container.NewVBox(
		headerContainer,
		widget.NewSeparator(),
		contentBox,
		widget.NewSeparator(),
		logsContainer,
	)

	return screen
}
