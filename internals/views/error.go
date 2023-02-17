package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/lenforiee/PassboltGUI/internals/controllers"
)

type ErrorWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

func NewErrorWindow(app *controllers.AppContext, err string) *ErrorWindow {

	window := (*app.App).NewWindow("PassboltGUI Error")
	view := &ErrorWindow{
		Window:    &window,
		Container: nil,
	}
	errorLabel := widget.NewLabel(err)
	button := widget.NewButton("OK", func() {
		(*view.Window).Close()
	})

	containerBox := fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		errorLabel,
		button,
	)
	view.Container = containerBox

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(fyne.NewSize(200, 100))
	(*view.Window).CenterOnScreen()
	return view
}
