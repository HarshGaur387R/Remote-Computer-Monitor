package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed images/loading_image.png
var embeddedImage embed.FS

func main() {
	myApp := app.New()

	myWindow := myApp.NewWindow("RCM")
	myWindow.Resize(fyne.NewSize(600, 600))
	myWindow.SetFixedSize(true)

	// Load embedded image
	imageData, err := embeddedImage.ReadFile("images/loading_image.png")
	if err != nil {
		panic(err)
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		panic(err)
	}

	embeddedImg := canvas.NewImageFromReader(bytes.NewReader(imageData), "loading_image."+format)
	embeddedImg.FillMode = canvas.ImageFillCover

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"Home",
			theme.HomeIcon(),
			HomeTab(myWindow)),
		container.NewTabItemWithIcon(
			"Agent",
			theme.AccountIcon(),
			AgentTab(myWindow)),
		container.NewTabItemWithIcon(
			"Connect",
			theme.ComputerIcon(),
			widget.NewLabel("Hello from connect page")),
		container.NewTabItemWithIcon(
			"About",
			theme.HelpIcon(),
			widget.NewLabel("Hello from HelpIcon")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(embeddedImg)

	go func() {
		time.Sleep(3 * time.Second)
		fyne.Do(func() {
			myWindow.SetContent(tabs)
		})
	}()

	myWindow.ShowAndRun()
}
