package main

import (
	// "bufio"
	"bufio"
	"encoding/json"
	"fmt"
	agentapis "guiinstaller/internal/agent-apis"
	"guiinstaller/internal/components"
	"guiinstaller/internal/constants"
	"guiinstaller/internal/utils"
	"image/color"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/fsnotify/fsnotify"
)

const logFilePath = constants.LOG_FILE_PATH

// LogEntry mirrors the JSON structure written by the agent logger.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// watchLogsAndDisplay tails logFilePath, parses each new JSON line, and
// appends a formatted string to logState. It blocks until cancel is called.
func watchLogsAndDisplay(entry *components.ReadOnlyEntry, cancel <-chan struct{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fyne.Do(func() { entry.Append(fmt.Sprintf("[watcher error] %v\n", err)) })
		return
	}
	defer watcher.Close()

	f, err := os.Open(logFilePath)
	if err != nil {
		fyne.Do(func() { entry.Append(fmt.Sprintf("[log open error] %v\n", err)) })
		return
	}
	defer f.Close()

	f.Seek(0, io.SeekStart)
	reader := bufio.NewReader(f)

	if err := watcher.Add(logFilePath); err != nil {
		fyne.Do(func() { entry.Append(fmt.Sprintf("[watcher add error] %v\n", err)) })
		return
	}

	appendLine := func(line string) {
		line = strings.TrimSpace(line)
		if line == "" {
			return
		}

		var entry_data struct {
			Timestamp string `json:"timestamp"`
			Level     string `json:"level"`
			Message   string `json:"message"`
		}

		var formatted string
		if err := json.Unmarshal([]byte(line), &entry_data); err != nil {
			// Show the parse error alongside the raw line so you can diagnose it
			formatted = fmt.Sprintf("[parse error: %v] raw: %s\n", err, line)
		} else {
			formatted = fmt.Sprintf("[%s] %s: %s\n", entry_data.Timestamp, entry_data.Level, entry_data.Message)
		}

		fyne.Do(func() { entry.Append(formatted) })
	}

	readNewLines := func() {
		for {
			line, err := reader.ReadString('\n')
			if len(strings.TrimSpace(line)) > 0 {
				appendLine(line)
			}
			if err != nil {
				if err != io.EOF {
					fyne.Do(func() { entry.Append(fmt.Sprintf("[read error] %v\n", err)) })
				}
				break
			}
		}
	}

	readNewLines()

	for {
		select {
		case <-cancel:
			return

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				readNewLines()
			}
			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				f.Close()
				f, err = os.Open(logFilePath)
				if err != nil {
					fyne.Do(func() { entry.Append(fmt.Sprintf("[reopen error] %v\n", err)) })
					return
				}
				reader.Reset(f)
				watcher.Add(logFilePath)
			}

		case watchErr, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fyne.Do(func() { entry.Append(fmt.Sprintf("[watcher error] %v\n", watchErr)) })
		}
	}
}

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

func AgentTab(parent fyne.Window) fyne.CanvasObject {
	// ==================== STATES ====================
	logState := binding.NewStringList()

	logsReadOnlyEntry := components.NewReadOnlyMultiLineEntry()
	logsReadOnlyEntry.SetMinRowsVisible(20)

	cancelWatch := make(chan struct{})
	go watchLogsAndDisplay(logsReadOnlyEntry, cancelWatch)

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

	logs, _ := logState.Get()

	for i := range logs {
		logsReadOnlyEntry.Append(logs[i])
	}

	// Combine background and logs
	logsContent := container.NewStack(
		container.NewVBox(
			logsReadOnlyEntry,
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

	parent.SetOnClosed(func() {
		close(cancelWatch)
	})

	return screen
}
