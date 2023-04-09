package views

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/internals/settings"
	"github.com/lenforiee/AmnesiaGUI/app/usecases/passbolt"
	"github.com/lenforiee/AmnesiaGUI/bundles"
	"github.com/sqweek/dialog"
)

type LoginView struct {
	Window    fyne.Window
	Container *fyne.Container
}

func NewLoginView(ctx app.AppContext) LoginView {

	window := ctx.App.NewWindow(fmt.Sprintf("%s :: Login", ctx.AppName))
	view := LoginView{
		Window: window,
	}

	userAgentLabel := widget.NewLabelWithStyle(
		"Passbolt User Agent",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemUserAgent := widget.NewEntry()
	itemUserAgent.SetText(ctx.UserSettings.UserAgent)

	checkUserAgent := ctx.UserSettings.UserAgent != ""
	userAgentEnable := widget.NewCheckWithData("Use Custom User Agent", binding.BindBool(&checkUserAgent))
	userAgentEnable.OnChanged = func(checked bool) {
		itemUserAgent.Disable()
		if checked {
			itemUserAgent.Enable()
		}
	}

	serverURILabel := widget.NewLabelWithStyle(
		"Passbolt Server URI",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemServerURI := widget.NewEntry()
	itemServerURI.SetText(ctx.UserSettings.ServerURI)

	privateKeyPathLabel := widget.NewLabelWithStyle(
		"Passbolt Private Key Path",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemPrivateKeyPath := widget.NewEntry()
	itemPrivateKeyPath.SetText(ctx.UserSettings.PrivateKeyPath)

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

	rememberInfo := widget.NewCheckWithData("Remember Info", binding.BindBool(&ctx.UserSettings.RememberMe))
	rememberInfo.OnChanged = func(checked bool) {
		if checked {
			ctx.UserSettings = settings.UserSettings{
				ServerURI:      itemServerURI.Text,
				PrivateKeyPath: itemPrivateKeyPath.Text,
				RememberMe:     checked,
			}

			if err := ctx.UserSettings.SaveUserSettings(); err != nil {
				errMsg := fmt.Sprintf("There was error while saving user settings: %s", err)
				logger.LogErr.Print(errMsg)

				errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
				errView.Window.Show()
				rememberInfo.SetChecked(false)
				return
			}

			itemServerURI.Disable()
			itemPrivateKeyPath.Disable()
			dialogBtn.Disable()
			itemUserAgent.Disable()
			userAgentEnable.Disable()

			itemPasswd.FocusGained()
		} else {
			itemServerURI.Enable()
			itemPrivateKeyPath.Enable()
			dialogBtn.Enable()
			userAgentEnable.Enable()

			if userAgentEnable.Checked {
				itemUserAgent.Enable()
			}

		}
	}

	loginFunc := func(password string) {
		go OnClickLogin(ctx, password, window)
	}

	loginButton := widget.NewButton("Login", func() {
		loginFunc(itemPasswd.Text)
	})

	itemPasswd.OnSubmitted = func(_ string) {
		loginFunc(itemPasswd.Text)
	}

	image := canvas.NewImageFromResource(bundles.ResourceAmnesiaLogoPng)
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
			userAgentLabel,
			itemUserAgent,
			userAgentEnable,
			widget.NewSeparator(),
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

	view.Window = window
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(350, 100))
	view.Window.CenterOnScreen()
	return view
}

func OnClickLogin(ctx app.AppContext, password string, loginWindow fyne.Window) {
	loadingSplash := NewLoadingSplash(ctx, "Logging in...")
	if err := passbolt.InitialisePassboltConnector(ctx, password); err != nil {
		errMsg := fmt.Sprintf("There was error while initialising passbolt client: %s", err)
		logger.LogErr.Println(errMsg)

		// Passbolt connector is bit weird, so we need to do this to get proper error message.
		errProperMessage := "Unknown error occured while initialising passbolt client."
		switch {
		case strings.Contains(err.Error(), "private key checksum failure"):
			errProperMessage = "Invalid passphrase. Please try again."
		case strings.Contains(err.Error(), "no such host"):
			errProperMessage = "Could not resolve server host. Please check your server url."
		case strings.Contains(err.Error(), "connection refused"):
			errProperMessage = "Could not connect to server. Please check your internet connection or server url."
		}

		errView := NewErrorView(ctx.App, ctx.AppName, errProperMessage, false)
		errView.Window.Show()
		return
	}

	loginWindow.Close()
	loadingSplash.Window.Close()

	// view, size := NewListWindow(app)
	// app.UpdateMainWindow(view.Window, size)
}
