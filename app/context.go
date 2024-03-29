package app

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/lenforiee/Amnesia/app/internals/logger"
	"github.com/lenforiee/Amnesia/app/internals/settings"
	"github.com/lenforiee/Amnesia/bundles"
	"github.com/passbolt/go-passbolt/api"
)

type AppContext struct {
	AppName        string
	MainWindow     fyne.Window
	App            fyne.App
	Context        context.Context
	PassboltClient *api.Client
	UserSettings   settings.UserSettings
}

func NewAppContext() AppContext {
	return AppContext{}
}

func (a *AppContext) InitialiseSystemTray() {
	desk, ok := a.App.(desktop.App)
	if ok {
		item := fyne.NewMenuItem("Show", func() {
			logger.LogInfo.Println("Showing main window")
			a.MainWindow.Show()
		})

		m := fyne.NewMenu(a.AppName, item)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(bundles.ResourceLogoPng)
	}
}

func (a *AppContext) UpdateMainWindow(window fyne.Window, size fyne.Size, centre bool) {

	logger.LogInfo.Println("Updating main window")
	a.MainWindow.SetTitle(window.Title())
	a.MainWindow.SetContent(window.Content())
	a.MainWindow.Resize(size)
	if centre {
		a.MainWindow.CenterOnScreen()
	}
}

func (a *AppContext) UpdateView(title string, content fyne.CanvasObject) {

	logger.LogInfo.Println("Updating view")
	a.MainWindow.SetTitle(title)
	a.MainWindow.SetContent(content)
}

func (a *AppContext) SetAppName(name string) {
	a.AppName = name
}

func (a *AppContext) SetMainWindow(window fyne.Window) {
	a.MainWindow = window
	a.MainWindow.SetMaster()
}

func (a *AppContext) SetApp(app fyne.App) {
	a.App = app
}

func (a *AppContext) SetContext(ctx context.Context) {
	a.Context = ctx
}

func (a *AppContext) SetPassboltClient(client *api.Client) {
	a.PassboltClient = client
}

func (a *AppContext) SetUserSettings(settings settings.UserSettings) {
	a.UserSettings = settings
}
