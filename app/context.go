package app

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/lenforiee/AmnesiaGUI/app/internals/settings"
	"github.com/lenforiee/AmnesiaGUI/bundles"
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
	if desk, ok := a.App.(desktop.App); ok {
		item := fyne.NewMenuItem("Show", func() {
			a.MainWindow.Show()
			// TODO: Add a gourotine to check if the connector is still connected.
		})

		m := fyne.NewMenu(a.AppName, item)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(bundles.ResourceLogoPng)
	}
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
