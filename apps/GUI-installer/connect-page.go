package main

import (
	"fmt"
	"guiinstaller/internal/constants"
	"guiinstaller/internal/utils"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func getContentfromConfig() (string, error) {

	bytes, err := os.ReadFile(constants.CONFIG_FILE_PATH)

	if err != nil {
		return "", err
	}

	content := string(bytes)
	return content, nil
}

func ConnectTab(parent fyne.Window) fyne.CanvasObject {

	// States
	agentFileExistState := binding.NewBool()
	configFileExistState := binding.NewBool()
	showErrorContainerState := binding.NewBool()
	showLoadingContinerState := binding.NewBool()

	showLoadingContinerState.Set(true)
	showErrorContainerState.Set(false)

	agentFileExistState.Set(utils.FileExists(constants.RCMA_BINARY_PATH))
	configFileExistState.Set(utils.FileExists(constants.CONFIG_FILE_PATH))

	header := canvas.NewText("finding", color.Opaque)
	loadingProgressBar := widget.NewProgressBarInfinite()
	loadingProgressBar.Hide()
	header.TextSize = 24
	message := widget.NewLabel("")

	// 1. Setup our target Fyne canvas image frame
	qrCanvasImage := canvas.NewImageFromImage(nil)
	qrCanvasImage.FillMode = canvas.ImageFillContain
	qrCanvasImage.SetMinSize(fyne.NewSize(256, 256))

	setErrorContainer := func(headertext string, messagetext string) {
		header.Text = headertext
		header.Alignment = fyne.TextAlignCenter
		header.Color = color.RGBA{R: 200, B: 0, G: 0, A: 255}
		message.Text = messagetext
	}

	GetErrorContainer := func() *fyne.Container {
		return container.NewVBox(header, message)
	}

	GetLoadingContainer := func() *fyne.Container {
		loadingHeader := canvas.NewText("Loading", color.Opaque)
		loadingHeader.Alignment = fyne.TextAlignCenter
		loadingHeader.Color = color.RGBA{R: 0, G: 0, B: 200, A: 255}
		loadingHeader.TextSize = 24
		loadingHeader.TextStyle.Bold = true

		return container.NewCenter(container.NewVBox(loadingHeader, loadingProgressBar))
	}

	GetMainContainer := func() *fyne.Container {
		content, err := getContentfromConfig()

		if err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			qrErrormessage := canvas.NewText("Error on gathering data for qr code", color.RGBA{R: 200, A: 255})
			qrErrormessage.Alignment = fyne.TextAlignCenter
			return container.NewCenter(container.NewVBox(header, message, qrErrormessage))
		}

		if content == "" {
			fmt.Println("config file has no content, Please contact dev team")
			qrErrormessage := canvas.NewText(
				"Config file is empty, please contect to dev team.",
				color.RGBA{R: 200, A: 255},
			)
			qrErrormessage.Alignment = fyne.TextAlignCenter
			return container.NewCenter(container.NewVBox(header, message, qrErrormessage))
		}

		img, err := utils.GenerateQrCode(content)
		if err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			qrErrormessage := canvas.NewText(
				"Error on generating QR code, please contact to dev team.",
				color.RGBA{R: 200, A: 255},
			)
			qrErrormessage.Alignment = fyne.TextAlignCenter
			return container.NewCenter(container.NewVBox(header, message, qrErrormessage))
		}

		// Update the existing canvas pointer with the new image asset
		qrCanvasImage.Image = img
		refreshBtn := widget.NewButton("Refresh", func() { fmt.Println("Hello world") })

		return container.NewCenter(container.NewVBox(header, message, qrCanvasImage, widget.NewSeparator(), widget.NewSeparator(), refreshBtn))
	}

	// Checking if agent binary and its config file exist
	isAgentFileExist, err := agentFileExistState.Get()
	showLoadingContinerState.Set(false) // Set loading to false before moving to checks

	if err != nil {
		setErrorContainer("Error on finding agent", "Facing error on finding agent binary please contact to dev team.")

	} else {
		if !isAgentFileExist {
			setErrorContainer("Agent Not Found", "Install the agent first, If its already installed then reinstall it.")
		} else {
			isConfigFileExist, err := configFileExistState.Get()
			if err != nil {
				setErrorContainer(
					"Error on finding config.json",
					"Start the agent wait for its activation then comback here.")
			} else {
				if !isConfigFileExist {
					setErrorContainer(
						"Config.json does not exist",
						"Start the agent wait for its activation, or contact dev.")

				} else {
					// Show Qr Code, Connected device, button to refresh, button to change auth
					header.Text = "Agent found, config found"
					header.Alignment = fyne.TextAlignCenter
					message.Text = "Agent binary and config file is found now you can connect"
				}
			}

		}
	}

	errorContainerState, err := showErrorContainerState.Get()
	if err != nil {
		fmt.Println("Error on reading showErrorContainerState")
		header := canvas.NewText("Error on reading showLoadingContainerState", nil)
		header.Alignment = fyne.TextAlignCenter
		header.Color = color.RGBA{R: 200, B: 0, G: 0, A: 255}
		header.TextSize = 20
		screen := container.NewCenter(
			container.NewVBox(
				header,
				widget.NewLabel("Error on reading showErrorContainerState, Pease contact to dev team."),
			),
		)
		return screen
	}

	loadingContainerState, err := showLoadingContinerState.Get()
	if err != nil {
		fmt.Println("Error on reading showLoadingContainerState")
		header := canvas.NewText("Error on reading showLoadingContainerState", nil)
		header.Alignment = fyne.TextAlignCenter
		header.Color = color.RGBA{R: 200, B: 0, G: 0, A: 255}
		header.TextSize = 20
		screen := container.NewCenter(
			container.NewVBox(
				header,
				widget.NewLabel("Error on reading showLoadingContainerState, Pease contact to dev team."),
			),
		)
		return screen
	}

	if errorContainerState {
		screen := GetErrorContainer()
		return screen
	} else if loadingContainerState {
		screen := GetLoadingContainer()
		loadingProgressBar.Show()
		return screen
	} else {
		screen := GetMainContainer()
		return screen
	}
}
