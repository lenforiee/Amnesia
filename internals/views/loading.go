package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/AmnesiaGUI/bundle"
	"github.com/lenforiee/AmnesiaGUI/internals/contexts"
)

type LoadingWindow struct {
	Window    *fyne.Window
	Label     *widget.Label
	Container *fyne.Container
}

func NewLoadingWindow(app *contexts.AppContext, text string) *LoadingWindow {
	var window fyne.Window
	if drv, ok := (*app.App).Driver().(desktop.Driver); ok {
		window = drv.CreateSplashWindow()
		window.SetTitle(fmt.Sprintf("%s :: Loading", app.AppName))
	}

	view := &LoadingWindow{
		Window:    nil,
		Container: nil,
	}

	loadingLabel := widget.NewLabelWithStyle(
		text,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	image := canvas.NewImageFromResource(bundle.ResourceAssetsImagesAmnesialogoPng)
	image.FillMode = canvas.ImageFillOriginal

	containerBox := container.NewBorder(
		nil,
		loadingLabel,
		nil,
		nil,
		image,
	)

	view.Window = &window
	view.Label = loadingLabel
	view.Container = containerBox

	size := fyne.NewSize(300, 200)

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(size)
	(*view.Window).CenterOnScreen()
	return view
}

func (view *LoadingWindow) UpdateText(text string) {
	(*view.Label).SetText(text)
}

func (view *LoadingWindow) StopLoading() {
	(*view.Window).Close()
}
