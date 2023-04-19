package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/Amnesia/app/internals/logger"
)

type ErrorView struct {
	Window    fyne.Window
	Container *fyne.Container
}

func NewErrorView(app fyne.App, appName string, err string, crash bool) ErrorView {

	logger.LogInfo.Printf("Creating new error view with message: %s, crash: %t", err, crash)
	window := app.NewWindow(fmt.Sprintf("%s :: Error", appName))
	view := ErrorView{
		Window: window,
	}

	labelHeader := "Error has occured!"
	if crash {
		labelHeader = "Amnesia has crashed!"
	}

	errorLabel := widget.NewLabelWithStyle(labelHeader,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	errorInfo := widget.NewLabel(err)
	errorInfo.Alignment = fyne.TextAlignCenter
	errorInfo.Wrapping = fyne.TextWrapWord

	button := widget.NewButton("OK", func() {
		view.Window.Close()
	})

	containerBox := container.NewBorder(
		errorLabel,
		button,
		nil,
		nil,
		errorInfo,
	)
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(400, 100))
	view.Window.CenterOnScreen()
	return view
}
