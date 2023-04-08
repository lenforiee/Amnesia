package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/AmnesiaGUI/bundle"
	"github.com/lenforiee/AmnesiaGUI/internals/contexts"
	"github.com/lenforiee/AmnesiaGUI/utils"
	"github.com/lenforiee/AmnesiaGUI/utils/logger"
	"github.com/sqweek/dialog"
)

type LoginWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

func NewLoginWindow(app *contexts.AppContext) (*LoginWindow, fyne.Size) {

	window := (*app.App).NewWindow(fmt.Sprintf("%s :: Login", app.AppName))
	view := &LoginWindow{
		Window:    &window,
		Container: nil,
	}

	serverURILabel := widget.NewLabelWithStyle(
		"Passbolt Server URI",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemServerURI := widget.NewEntry()
	itemServerURI.SetText(app.UserConfig.ServerURI)

	privateKeyPathLabel := widget.NewLabelWithStyle(
		"Passbolt Private Key Path",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemPrivateKeyPath := widget.NewEntry()
	itemPrivateKeyPath.SetText(app.UserConfig.PrivateKeyPath)

	dialogBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		filename, err := dialog.File().Filter("Passbolt Private Key File (.txt, .pem)", "txt", "pem").Load()
		if err != nil {
			return
		}
		itemPrivateKeyPath.SetText(filename)
	})

	passwdLabel := widget.NewLabelWithStyle(
		"Passbolt Passphrase",
		fyne.TextAlignCenter,
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
				app.CreateNewWindowAndShow(errView.Window)
				rememberInfo.SetChecked(false)
				return
			}
			itemServerURI.Disable()
			itemPrivateKeyPath.Disable()
			dialogBtn.Disable()

			// have focus on password item.
			itemPasswd.FocusGained()
		} else {
			itemServerURI.Enable()
			itemPrivateKeyPath.Enable()
			dialogBtn.Enable()

			itemPasswd.FocusLost()
		}
	}

	loginButton := widget.NewButton("Login", func() {
		loadingSplash := NewLoadingWindow(app, "Logging in...")
		app.CreateNewWindowAndShow(loadingSplash.Window)
		OnClickLogin(app, itemPasswd.Text)
		loadingSplash.StopLoading()
	})

	image := canvas.NewImageFromResource(bundle.ResourceAssetsImagesAmnesialogoPng)
	image.FillMode = canvas.ImageFillOriginal

	containerBox := container.NewBorder(
		image,
		container.New(
			layout.NewVBoxLayout(),
			loginButton,
		),
		nil,
		nil,
		container.New(
			layout.NewVBoxLayout(),
			serverURILabel,
			itemServerURI,
			privateKeyPathLabel,
			container.NewBorder(
				nil,
				nil,
				nil,
				dialogBtn,
				itemPrivateKeyPath,
			),
			passwdLabel,
			itemPasswd,
			rememberInfo,
		),
	)

	view.Window = &window
	view.Container = containerBox

	size := fyne.NewSize(350, 100)

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(size)
	(*view.Window).CenterOnScreen()
	return view, size
}

func OnClickLogin(app *contexts.AppContext, password string) {
	if err := app.InitialisePassbolt(password); err != nil {
		errMsg := fmt.Sprintf("There was error while initialising passbolt client: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorWindow(app, errMsg)
		app.CreateNewWindowAndShow(errView.Window)
		return
	}

	view, size := NewListWindow(app)
	app.UpdateMainWindow(view.Window, size)
}
