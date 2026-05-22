package main

import (
	"encoding/json"
	"fmt"
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

type Config struct {
	LANIP     string `json:"lan_ip"`
	Port      int    `json:"port"`
	AuthToken string `json:"auth_token"`
}

func getContentFromConfig() (string, error) {

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
	showLoadingContainerState := binding.NewBool()

	showLoadingContainerState.Set(true)
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

	RefreshQrCode := func() {
		content, err := getContentFromConfig()

		if err != nil {
			fmt.Println(fmt.Errorf("Error on refresh:  %v", err))
			dialog.ShowError(fmt.Errorf("%v", err), parent)
			return
		}

		if content == "" {
			fmt.Println("config file has no content, Please contact dev team")
			dialog.ShowError(fmt.Errorf("config file has no content, Please contact dev team"), parent)
			return
		}

		img, err := utils.GenerateQrCode(content)
		if err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			dialog.ShowError(fmt.Errorf("Error on generating qr-code: %v", err), parent)
			return
		}

		if img == nil {
			fmt.Println(fmt.Errorf("Error on refresh Generated image is empty, please contact to dev team"))
			dialog.ShowError(fmt.Errorf("Error on generating qr-code: Image is nil"), parent)
			return
		}

		// Update the existing canvas pointer with the new image asset
		qrCanvasImage.Image = img
		qrCanvasImage.Refresh()
	}

	ResetAuth := func() {
		content, err := getContentFromConfig()
		if err != nil || content == "" {
			fmt.Println(fmt.Errorf("Error on reset auth:  %v", err))
			dialog.ShowError(fmt.Errorf("%v", err), parent)
			return
		}

		var config Config

		json_err := json.Unmarshal([]byte(content), &config)
		if json_err != nil {
			fmt.Println(fmt.Errorf("Error unmarshaling JSON: %v", json_err))
			dialog.ShowError(fmt.Errorf("%v", json_err), parent)
			return
		}

		// Now you can access the fields normally
		fmt.Printf("LANIP: %s, Port: %d, AuthToken: %s\n", config.LANIP, config.Port, config.AuthToken)

		token, err := utils.GenerateSecureToken(32)
		if err != nil {
			fmt.Println(fmt.Errorf("Error on generating new secure token: %v", err))
			dialog.ShowError(fmt.Errorf("%v", err), parent)
			return
		}

		config.AuthToken = token

		new_json_data, json_mar_err := json.Marshal(config)
		if json_mar_err != nil {
			fmt.Println(fmt.Errorf("Error marshaling JSON: %v", json_mar_err))
			dialog.ShowError(fmt.Errorf("%v", json_mar_err), parent)
			return
		}

		new_json_data_string := string(new_json_data)
		new_json_data_bytes := []byte(new_json_data_string)

		write_file_err := os.WriteFile(constants.CONFIG_FILE_PATH, new_json_data_bytes, 0644)
		if write_file_err != nil {
			fmt.Println(fmt.Errorf("Error on writing data to config.json: %v", write_file_err))
			dialog.ShowError(fmt.Errorf("%v", write_file_err), parent)
			return
		}

		img, err := utils.GenerateQrCode(new_json_data_string)
		if err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			dialog.ShowError(fmt.Errorf("Error on generating qr-code: %v", err), parent)
			return
		}

		if img == nil {
			fmt.Println(fmt.Errorf("Generated image is empty, please contact to dev team"))
			dialog.ShowError(fmt.Errorf("Error on generating qr-code: image is nil"), parent)
			return
		}

		// Update the existing canvas pointer with the new image asset
		qrCanvasImage.Image = img
		qrCanvasImage.Refresh()
	}

	setErrorContainer := func(headerText string, messageText string) {
		header.Text = headerText
		header.Alignment = fyne.TextAlignCenter
		header.Color = color.RGBA{R: 200, B: 0, G: 0, A: 255}
		message.Text = messageText

		showErrorContainerState.Set(true)
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
		content, err := getContentFromConfig()

		if err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			qrErrormessage := canvas.NewText("Error on gathering data for qr code", color.RGBA{R: 200, A: 255})
			qrErrormessage.Alignment = fyne.TextAlignCenter
			return container.NewCenter(container.NewVBox(header, message, qrErrormessage))
		}

		if content == "" {
			fmt.Println("config file has no content, Please contact dev team")
			qrErrormessage := canvas.NewText(
				"Config file is empty, please contact to dev team.",
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

		if img == nil {
			fmt.Println(fmt.Errorf("Failed to generate QRCode, image is empty"))
			qrErrormessage := canvas.NewText(
				"Error on generating QR code. Image is empty, please contact to dev team.",
				color.RGBA{R: 200, A: 255},
			)
			qrErrormessage.Alignment = fyne.TextAlignCenter
			return container.NewCenter(container.NewVBox(header, message, qrErrormessage))
		}

		// Update the existing canvas pointer with the new image asset
		qrCanvasImage.Image = img
		refreshBtn := widget.NewButton("Refresh", func() { RefreshQrCode() })
		resetAuthToken := widget.NewButton("Reset auth token", func() { ResetAuth() })

		return container.NewCenter(container.NewVBox(
			header,
			message,
			qrCanvasImage,
			widget.NewSeparator(),
			widget.NewSeparator(),
			refreshBtn,
			widget.NewSeparator(),
			resetAuthToken,
		))
	}

	// Checking if agent binary and its config file exist
	isAgentFileExist, err := agentFileExistState.Get()
	showLoadingContainerState.Set(false) // Set loading to false before moving to checks

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
					"Start the agent wait for its activation then comeback here.")
			} else {
				if !isConfigFileExist {
					setErrorContainer(
						"Config.json does not exist",
						"Start the agent wait for its activation, or contact dev.")

				} else {
					// Show Qr Code, Connected device, button to refresh, button to change auth
					header.Text = "Scan QR Code"
					header.Alignment = fyne.TextAlignCenter
					message.Text = "To connect with agent first start the agent \nthen scan this QR code using RCM mobile app."
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

	loadingContainerState, err := showLoadingContainerState.Get()
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
