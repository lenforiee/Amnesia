package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/bundles"
)

type LoadingSplash struct {
	Window    fyne.Window
	Label     widget.Label
	Container *fyne.Container
}

func NewLoadingSplash(ctx app.AppContext, text string) LoadingSplash {

	var window fyne.Window
	if drv, ok := ctx.App.Driver().(desktop.Driver); ok {
		window = drv.CreateSplashWindow()
		window.SetTitle(fmt.Sprintf("%s :: Loading...", ctx.AppName))
	}

	view := LoadingSplash{
		Window: window,
	}

	loadingLabel := widget.NewLabelWithStyle(
		text,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	image := canvas.NewImageFromResource(bundles.ResourceAmnesiaLogoPng)
	image.FillMode = canvas.ImageFillOriginal

	containerBox := container.NewBorder(
		nil,
		loadingLabel,
		nil,
		nil,
		image,
	)

	view.Label = *loadingLabel
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(300, 200))
	view.Window.CenterOnScreen()
	return view
}

func (s *LoadingSplash) UpdateText(text string) {
	s.Label.SetText(text)
}

func (s *LoadingSplash) Close() {
	s.Window.Close()
}
