package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/PassboltGUI/internals/controllers"
	"github.com/lenforiee/PassboltGUI/utils"
	"github.com/lenforiee/PassboltGUI/utils/logger"
	"github.com/sqweek/dialog"
)

type LoginWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

func NewLoginWindow(app *controllers.AppContext) (*LoginWindow, fyne.Size) {

	window := (*app.App).NewWindow("PassboltGUI Login")
	view := &LoginWindow{
		Window:    &window,
		Container: nil,
	}

	serverURILabel := widget.NewLabelWithStyle(
		"Enter your server URI",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	itemServerURI := widget.NewEntry()
	itemServerURI.SetText(app.UserConfig.ServerURI)

	privateKeyPathLabel := widget.NewLabelWithStyle(
		"Enter your private key path",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	itemPrivateKeyPath := widget.NewEntry()
	itemPrivateKeyPath.SetText(app.UserConfig.PrivateKeyPath)

	dialogBtn := widget.NewButton("Choose file", func() {
		filename, err := dialog.File().Filter("Passbolt Private Key File (.txt, .pem)", "txt", "pem").Load()
		if err != nil {
			return
		}
		itemPrivateKeyPath.SetText(filename)
	})

	passwdLabel := widget.NewLabelWithStyle(
		"Enter your passphrase",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	itemPasswd := widget.NewPasswordEntry()

	rememberInfo := widget.NewCheckWithData("Remember Info", binding.BindBool(&app.UserConfig.RememberMe))
	rememberInfo.OnChanged = func(checked bool) {
		if checked {
			app.UserConfig = &utils.UserConfig{
				ServerURI:      itemServerURI.Text,
				PrivateKeyPath: itemPrivateKeyPath.Text,
				RememberMe:     checked,
			}
			if err := app.SaveConfig(); err != nil {
				errMsg := fmt.Sprintf("There was error while saving user settings: %s", err)
				logger.LogErr.Println(errMsg)

				errView := NewErrorWindow(app, errMsg)
				app.CreateNewWindowWithView(errView.Window)
				rememberInfo.SetChecked(false)
				return
			}
			itemServerURI.Disable()
			itemPrivateKeyPath.Disable()
			dialogBtn.Disable()
		} else {
			itemServerURI.Enable()
			itemPrivateKeyPath.Enable()
			dialogBtn.Enable()
		}
	}

	loginButton := widget.NewButton("Login", func() {
		OnClickLogin(app, itemPasswd.Text)
		return
	})

	image := canvas.NewImageFromFile("./assets/logo_white.png")
	image.FillMode = canvas.ImageFillOriginal

	containerBox := container.New(
		layout.NewVBoxLayout(),
		image,
		widget.NewSeparator(),
		serverURILabel,
		itemServerURI,
		privateKeyPathLabel,
		container.New(
			layout.NewVBoxLayout(),
			itemPrivateKeyPath,
			dialogBtn,
		),
		passwdLabel,
		itemPasswd,
		rememberInfo,
		loginButton,
	)
	view.Window = &window
	view.Container = containerBox

	size := fyne.NewSize(350, 100)

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(size)
	(*view.Window).CenterOnScreen()
	return view, size
}

func OnClickLogin(app *controllers.AppContext, password string) {
	if err := app.InitialisePassbolt(password); err != nil {
		errMsg := fmt.Sprintf("There was error while initialising passbolt client: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorWindow(app, errMsg)
		app.CreateNewWindowWithView(errView.Window)
		return
	}

	view, size := NewListWindow(app)
	app.UpdateMainWindow(view.Window, size)
}
