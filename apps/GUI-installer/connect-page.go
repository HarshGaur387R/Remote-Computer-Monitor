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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	TokenLength     = 32
	ConfigFilePerms = 0644
	ErrorRed        = 200
	InfoBlue        = 200
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

// UpdateQrImage generates and updates QR code image from content
func UpdateQrImage(content string, qrCanvasImage *canvas.Image, parent fyne.Window) error {
	if content == "" {
		err := fmt.Errorf("config file has no content")
		dialog.ShowError(err, parent)
		return err
	}

	img, err := utils.GenerateQrCode(content)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Error generating QR code: %v", err), parent)
		return err
	}

	if img == nil {
		err := fmt.Errorf("generated QR image is empty")
		dialog.ShowError(err, parent)
		return err
	}

	qrCanvasImage.Image = img
	qrCanvasImage.Refresh()
	return nil
}

// UpdateConfigAuthToken generates a new auth token and saves it to config file
func UpdateConfigAuthToken(parent fyne.Window) (string, error) {
	content, err := getContentFromConfig()
	if err != nil || content == "" {
		dialog.ShowError(fmt.Errorf("failed to read config: %v", err), parent)
		return "", err
	}

	var config Config
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		dialog.ShowError(fmt.Errorf("invalid config format: %v", err), parent)
		return "", err
	}

	// Validate config structure
	if config.LANIP == "" || config.Port == 0 || config.AuthToken == "" {
		err := fmt.Errorf("incomplete config: missing required fields")
		dialog.ShowError(err, parent)
		return "", err
	}

	token, err := utils.GenerateSecureToken(TokenLength)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to generate token: %v", err), parent)
		return "", err
	}

	config.AuthToken = token
	newJSON, err := json.Marshal(config)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to marshal config: %v", err), parent)
		return "", err
	}

	if err := os.WriteFile(constants.CONFIG_FILE_PATH, newJSON, ConfigFilePerms); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save config: %v", err), parent)
		return "", err
	}

	return string(newJSON), nil
}

// CreateErrorText creates a styled error text element
func CreateErrorText(text string) *canvas.Text {
	errText := canvas.NewText(text, color.RGBA{R: ErrorRed, A: 255})
	errText.Alignment = fyne.TextAlignCenter
	return errText
}

func ConnectTab(parent fyne.Window) fyne.CanvasObject {

	// UI Elements
	header := canvas.NewText("finding", color.Opaque)
	loadingProgressBar := widget.NewProgressBarInfinite()
	loadingProgressBar.Hide()
	header.TextSize = 24
	message := widget.NewLabel("")

	// State tracking
	var showError bool
	var showLoading bool

	// Setup QR code image canvas
	qrCanvasImage := canvas.NewImageFromImage(nil)
	qrCanvasImage.FillMode = canvas.ImageFillContain
	qrCanvasImage.SetMinSize(fyne.NewSize(256, 256))

	RefreshQrCode := func() {
		content, err := getContentFromConfig()
		if err != nil {
			fmt.Println(fmt.Errorf("Error on refresh: %v", err))
			dialog.ShowError(err, parent)
			return
		}

		if err := UpdateQrImage(content, qrCanvasImage, parent); err != nil {
			fmt.Println(fmt.Errorf("Error updating QR image: %v", err))
		}
	}

	ResetAuth := func() {
		newConfig, err := UpdateConfigAuthToken(parent)
		if err != nil {
			fmt.Println(fmt.Errorf("Error resetting auth: %v", err))
			return
		}

		// Update QR code with new token
		if err := UpdateQrImage(newConfig, qrCanvasImage, parent); err != nil {
			fmt.Println(fmt.Errorf("Error updating QR after reset: %v", err))
			return
		}

		dialog.ShowInformation("Success", "Auth token has been reset", parent)
	}

	setErrorContainer := func(headerText string, messageText string) {
		header.Text = headerText
		header.Alignment = fyne.TextAlignCenter
		header.Color = color.RGBA{R: ErrorRed, B: 0, G: 0, A: 255}
		message.Text = messageText
		showError = true
		showLoading = false
	}

	GetErrorContainer := func() *fyne.Container {
		return container.NewVBox(header, message)
	}

	GetLoadingContainer := func() *fyne.Container {
		loadingHeader := canvas.NewText("Loading", color.Opaque)
		loadingHeader.Alignment = fyne.TextAlignCenter
		loadingHeader.Color = color.RGBA{R: 0, G: 0, B: InfoBlue, A: 255}
		loadingHeader.TextSize = 24
		loadingHeader.TextStyle.Bold = true

		return container.NewCenter(container.NewVBox(loadingHeader, loadingProgressBar))
	}

	GetMainContainer := func() *fyne.Container {
		content, err := getContentFromConfig()

		if err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			return container.NewCenter(container.NewVBox(header, message, CreateErrorText("Error on gathering data for qr code")))
		}

		if content == "" {
			fmt.Println("config file has no content, Please contact dev team")
			return container.NewCenter(container.NewVBox(header, message, CreateErrorText("Config file is empty, please contact dev team.")))
		}

		if err := UpdateQrImage(content, qrCanvasImage, parent); err != nil {
			fmt.Println(fmt.Errorf("%v", err))
			return container.NewCenter(container.NewVBox(header, message, CreateErrorText("Error on generating QR code, please contact dev team.")))
		}

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

	// Check if agent binary and config file exist
	agentExists := utils.FileExists(constants.RCMA_BINARY_PATH)
	configExists := utils.FileExists(constants.CONFIG_FILE_PATH)

	if !agentExists {
		setErrorContainer("Agent Not Found", "Install the agent first, If its already installed then reinstall it.")
	} else if !configExists {
		setErrorContainer("Config.json does not exist", "Start the agent wait for its activation, or contact dev.")
	} else {
		// All checks passed, show main content
		header.Text = "Scan QR Code"
		header.Alignment = fyne.TextAlignCenter
		message.Text = "To connect with agent first start the agent \nthen scan this QR code using RCM mobile app."
		showError = false
		showLoading = false
	}

	// Render logic: Check error state first, then loading, then main content
	// This allows errors to take priority over normal UI
	if showError {
		return GetErrorContainer()
	} else if showLoading {
		loadingProgressBar.Show()
		return GetLoadingContainer()
	} else {
		return GetMainContainer()
	}
}
