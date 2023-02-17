package controllers

import (
	"context"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/lenforiee/AmnesiaGUI/utils"
	"github.com/passbolt/go-passbolt/api"
)

type AppContext struct {
	AppName        string
	MainWindow     *fyne.Window
	App            *fyne.App
	Context        *context.Context
	PassboltClient *api.Client
	UserConfig     *utils.UserConfig
}

func NewAppContext(appName string, app *fyne.App, window *fyne.Window, context *context.Context) *AppContext {
	return &AppContext{
		AppName:        appName,
		MainWindow:     window,
		App:            app,
		Context:        context,
		PassboltClient: nil,
		UserConfig:     nil,
	}
}

func (a *AppContext) InitialiseSystemTray(logo fyne.Resource) {
	if desk, ok := (*a.App).(desktop.App); ok {
		item := fyne.NewMenuItem("Show", func() {
			(*a.MainWindow).Show()
		})

		m := fyne.NewMenu(a.AppName, item)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(logo)
	}
}

func (a *AppContext) UpdateMainWindow(window *fyne.Window, size fyne.Size) {
	(*a.MainWindow).SetTitle((*window).Title())
	(*a.MainWindow).SetContent((*window).Content())
	(*a.MainWindow).Resize(size)
	(*a.MainWindow).CenterOnScreen()
}

func (a *AppContext) StartMainWindow() {
	(*a.MainWindow).ShowAndRun()
}

func (a *AppContext) CreateNewWindowAndShow(window *fyne.Window) {
	(*window).Show()
}

func (a *AppContext) LoadConfig() error {
	config, err := utils.LoadUserConfig()

	if err != nil {
		return err
	}
	a.UserConfig = config

	return nil
}

func (a *AppContext) SaveConfig() error {
	err := a.UserConfig.SaveUserConfig()
	if err != nil {
		return err
	}

	return nil
}

func (a *AppContext) InitialisePassbolt(password string) error {

	// read the private key file
	privateKey, err := ioutil.ReadFile(a.UserConfig.PrivateKeyPath)
	if err != nil {
		return err
	}

	client, err := api.NewClient(nil, "", a.UserConfig.ServerURI, string(privateKey), password)
	if err != nil {
		return err
	}

	err = client.Login(*a.Context)
	if err != nil {
		return err
	}

	a.PassboltClient = client
	return nil
}
