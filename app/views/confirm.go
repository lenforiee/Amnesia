package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
)

type ConfirmView struct {
	Window    fyne.Window
	OnYes     func()
	Container *fyne.Container
}

func NewConfirmView(ctx *amnesiaApp.AppContext, msg string) *ConfirmView {

	logger.LogInfo.Printf("Creating new confirm view with message: %s", msg)
	window := ctx.App.NewWindow(fmt.Sprintf("%s :: Confirm", ctx.AppName))
	view := &ConfirmView{
		Window: window,
	}
	errorLabel := widget.NewLabelWithStyle("Warning!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	errorInfo := widget.NewLabel(msg)
	errorInfo.Wrapping = fyne.TextWrapWord

	yesBtn := widget.NewButton("Yes", func() {
		view.Window.Close()
		view.OnYes()
	})

	noBtn := widget.NewButton("No", func() {
		view.Window.Close()
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

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(400, 100))
	view.Window.CenterOnScreen()
	return view
}

func (v *ConfirmView) SetOnYesEvent(callback func()) {
	v.OnYes = callback
}
