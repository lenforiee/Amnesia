package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"

	"github.com/lenforiee/PassboltGUI/internals/controllers"
	"github.com/lenforiee/PassboltGUI/models"
)

type ResourceWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

func NewResourceWindow(app *controllers.AppContext, resource *models.Resource) *ResourceWindow {

	window := (*app.App).NewWindow(fmt.Sprintf("PassboltGUI Resource: %s", resource.Name))
	view := &ResourceWindow{
		Window:    &window,
		Container: nil,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	itemName := widget.NewEntry()

	itemName.SetText(resource.Name)
	itemName.Disable()

	usernameLabel := widget.NewLabelWithStyle(
		"Username",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemUsername := widget.NewEntry()
	itemUsername.TextStyle = fyne.TextStyle{Bold: true}

	itemUsername.SetText(resource.Username)
	itemUsername.Disable()

	uriLabel := widget.NewLabelWithStyle(
		"URI",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemUri := widget.NewEntry()
	itemUri.TextStyle = fyne.TextStyle{Bold: true}

	itemUri.SetText(resource.URI)
	itemUri.Disable()

	passwdLabel := widget.NewLabelWithStyle(
		"Password",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemPasswd := widget.NewPasswordEntry()
	itemPasswd.TextStyle = fyne.TextStyle{Bold: true}

	itemPasswd.SetText(resource.Password)
	itemPasswd.Disable()

	descLabel := widget.NewLabelWithStyle(
		"Description",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemDesc := widget.NewEntry()
	itemDesc.TextStyle = fyne.TextStyle{Bold: true}

	itemDesc.SetText(resource.Description)
	itemDesc.Disable()

	copyUsername := widget.NewButton("Copy Username", func() {
		clipboard.Write(clipboard.FmtText, []byte(resource.Username))
	})

	copyPasswd := widget.NewButton("Copy Password", func() {
		clipboard.Write(clipboard.FmtText, []byte(resource.Password))
	})

	closeBtn := widget.NewButton("Close", func() {
		(*view.Window).Close()
	})

	containerBox := container.New(
		layout.NewVBoxLayout(),
		nameLabel,
		itemName,
		usernameLabel,
		itemUsername,
		uriLabel,
		itemUri,
		passwdLabel,
		itemPasswd,
		descLabel,
		itemDesc,
		copyUsername,
		copyPasswd,
		widget.NewSeparator(),
		closeBtn,
	)
	view.Container = containerBox

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(fyne.NewSize(350, 100))
	(*view.Window).CenterOnScreen()
	return view
}
