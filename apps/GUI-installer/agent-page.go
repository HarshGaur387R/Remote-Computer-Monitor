package main

import (
	"fmt"
	agentapis "guiinstaller/internal/agent-apis"
	"guiinstaller/internal/constants"
	"guiinstaller/internal/utils"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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

		// ==================== STATES =========================
		statusValue, _ := statusState.Get()
		contentBox.Objects = nil

		// ==================== CONTROL BAR ====================
		playBtn := widget.NewButton("▶ Start", func() {

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
		})
		pauseBtn := widget.NewButton("⏸ Stop", func() {
			// Pause button action
			statusState.Set(false)
		})
		restartBtn := widget.NewButton("🔄 Restart", func() {
			// Restart button action
		})

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
