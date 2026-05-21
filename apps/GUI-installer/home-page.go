package main

import (
	"fmt"
	"guiinstaller/internal/constants"
	"guiinstaller/internal/utils"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func handel_install(btn *widget.Button, installed binding.Bool, downloading binding.Bool, progressBar *widget.ProgressBar, progressStatusLabel *widget.Label, parent fyne.Window) {

	println("installing agent...")

	btn.Disable()
	println("Button disabled")

	// Show progress bar and reset it
	downloading.Set(true)
	progressBar.Value = 0

	go func() {

		fyne.Do(func() {
			btn.Enable()
			downloading.Set(false)
			println("Button enabled")
		})

		errOnDownloading := utils.DownloadAgent(constants.URL, constants.RCMA_BINARY_PATH, func(progress float64) {
			fyne.Do(func() {
				progressBar.Value = progress / 100.0 // Convert percentage to 0-1 range
				// Update progress status label
				progressStatusLabel.SetText(fmt.Sprintf("%.1f%%", progress))
			})
		})

		// Update UI back on main thread
		fyne.Do(func() {
			btn.Enable()
			downloading.Set(false)
			println("Button enabled")

			if errOnDownloading != nil {
				fmt.Printf("Error on downloading agent: %v\n", errOnDownloading)
				dialog.ShowError(errOnDownloading, parent)
				return
			}
			println("Agent downloaded!")
			installed.Set(true) // Update binding after successful installation
		})
	}()
}

func handel_uninstall(btn *widget.Button, installed binding.Bool, parent fyne.Window) {

	println("Uninstalling agent")

	// First check if agent is running or not. because running agent can not be deleted

	btn.Disable()
	println("Button disabled")

	dirs := []string{constants.RCMA_BINARY_DIR_PATH, constants.CONFIG_DIR}
	for _, dir := range dirs {
		err := utils.DeleteDir(dir)
		if err != nil {
			dialog.ShowError(err, parent)
		}
	}

	println("Agent removed")
	installed.Set(false)
}

func HomeTab(parent fyne.Window) fyne.CanvasObject {

	// STATES:
	downloading := binding.NewBool()
	downloading.Set(false)

	installed := binding.NewBool()
	installed.Set(utils.FileExists(constants.RCMA_BINARY_PATH))

	header := canvas.NewText("Welcome to RCM", color.RGBA{
		R: 0,
		G: 0,
		B: 200, // blue intensity
		A: 255, // fully opaque
	})
	header.TextStyle.Bold = true
	header.TextSize = 32

	headerContainer := container.NewCenter(header)

	// Create progress bar
	progressBar := widget.NewProgressBar()
	progressBar.Value = 0
	progressStatusLabel := widget.NewLabel("")

	// Create a container to hold the dynamic content
	contentBox := container.NewVBox()

	// Function to update UI based on installed state and downloading state
	updateUI := func() {
		value, _ := installed.Get()
		isDownloading, _ := downloading.Get()

		// Clear previous content
		contentBox.Objects = nil

		if isDownloading {
			// Show downloading state
			downloadingLabel := widget.NewLabel("Downloading agent...")
			contentBox.Add(downloadingLabel)
			contentBox.Add(progressBar)
			contentBox.Add(progressStatusLabel)
		} else {
			// Determine button label and message based on installed state
			var buttonLabel string
			var messageText string

			if !value {
				buttonLabel = "Install Agent"
				messageText = "Agent not found, Please install the agent to continue."
			} else {
				buttonLabel = "Uninstall Agent"
				messageText = "Agent Found! Check agent section for more info."
			}

			// Create button with dynamic handler
			button := widget.NewButton(buttonLabel, nil)
			button.OnTapped = func() {
				if !value {
					handel_install(button, installed, downloading, progressBar, progressStatusLabel, parent)
				} else {
					handel_uninstall(button, installed, parent)
				}
			}

			messageBox := widget.NewLabel(messageText)

			contentBox.Add(messageBox)
			contentBox.Add(container.NewPadded(button))
		}

		contentBox.Refresh()
	}

	// Set up listener to update UI when binding changes
	installed.AddListener(binding.NewDataListener(updateUI))
	downloading.AddListener(binding.NewDataListener(updateUI))

	// Call updateUI once to set initial state
	updateUI()

	screen := container.NewCenter(container.NewVBox(
		headerContainer,
		contentBox,
	))

	return screen
}
