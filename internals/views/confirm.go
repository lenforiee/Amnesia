package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/lenforiee/AmnesiaGUI/internals/contexts"
)

type ConfirmWindow struct {
	Window    *fyne.Window
	OnYes     func()
	Container *fyne.Container
}

func NewConfirmWindow(app *contexts.AppContext, msg string) *ConfirmWindow {

	window := (*app.App).NewWindow(fmt.Sprintf("%s :: Confirm", app.AppName))
	view := &ConfirmWindow{
		Window:    &window,
		Container: nil,
	}
	errorLabel := widget.NewLabelWithStyle("Warning!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	errorInfo := widget.NewLabel(msg)
	errorInfo.Wrapping = fyne.TextWrapWord

	yesBtn := widget.NewButton("Yes", func() {
		(*view.Window).Close()
		(*view).OnYes()
	})

	noBtn := widget.NewButton("No", func() {
		(*view.Window).Close()
	})

	containerBox := container.NewBorder(
		errorLabel,
		container.New(
			layout.NewGridLayout(2),
			yesBtn,
			noBtn,
		),
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
