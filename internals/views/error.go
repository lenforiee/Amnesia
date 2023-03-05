package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/lenforiee/AmnesiaGUI/internals/contexts"
)

type ErrorWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

func NewErrorWindow(app *contexts.AppContext, err string) *ErrorWindow {

	window := (*app.App).NewWindow(fmt.Sprintf("%s :: Error", app.AppName))
	view := &ErrorWindow{
		Window:    &window,
		Container: nil,
	}
	errorLabel := widget.NewLabelWithStyle("Error has occured!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	errorInfo := widget.NewLabel(err)
	errorInfo.Wrapping = fyne.TextWrapWord

	button := widget.NewButton("OK", func() {
		(*view.Window).Close()
	})

	containerBox := container.NewBorder(
		errorLabel,
		button,
		nil,
		nil,
		errorInfo,
	)
	view.Container = containerBox

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(fyne.NewSize(400, 100))
	(*view.Window).CenterOnScreen()
	return view
}
