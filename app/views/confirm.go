package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	amnesiaApp "github.com/lenforiee/Amnesia/app"
	"github.com/lenforiee/Amnesia/app/internals/logger"
)

type ConfirmView struct {
	Window    fyne.Window
	OnYes     func()
	OnNo      func()
	Container *fyne.Container
}

func NewConfirmView(ctx *amnesiaApp.AppContext, msg string) *ConfirmView {

	logger.LogInfo.Printf("Creating new confirm view with message: %s", msg)
	window := ctx.App.NewWindow(fmt.Sprintf("%s :: Confirm", ctx.AppName))
	view := &ConfirmView{
		Window: window,
	}
	warningLabel := widget.NewLabelWithStyle("Warning!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	warningInfo := widget.NewLabel(msg)
	warningInfo.Alignment = fyne.TextAlignCenter
	warningInfo.Wrapping = fyne.TextWrapWord

	yesBtn := widget.NewButton("Yes", func() {
		view.Window.Close()
		view.OnYes()
	})

	noBtn := widget.NewButton("No", func() {
		view.Window.Close()
		view.OnNo()
	})

	containerBox := container.NewBorder(
		warningLabel,
		container.New(
			layout.NewGridLayout(2),
			yesBtn,
			noBtn,
		),
		nil,
		nil,
		warningInfo,
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

func (v *ConfirmView) SetOnNoEvent(callback func()) {
	v.OnNo = callback
}
