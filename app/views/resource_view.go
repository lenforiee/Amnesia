package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/models"
)

type ResourceView struct {
	Window    fyne.Window
	Container *fyne.Container
}

func NewResourceView(ctx amnesiaApp.AppContext, token string, resource models.Resource) ResourceView {

	window := ctx.App.NewWindow(fmt.Sprintf("%s :: View Resource", ctx.AppName))

	view := ResourceView{
		Window: window,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemName := widget.NewEntry()
	itemName.SetPlaceHolder("eg. Amazon")

	itemName.SetText(resource.Name)
	itemName.Disable()

	usernameLabel := widget.NewLabelWithStyle(
		"Username",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemUsername := widget.NewEntry()
	itemUsername.SetPlaceHolder("eg. example@example.com")

	itemUsername.SetText(resource.Username)
	itemUsername.Disable()

	uriLabel := widget.NewLabelWithStyle(
		"URI",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemUri := widget.NewEntry()
	itemUri.SetPlaceHolder("eg. https://amazon.com")

	itemUri.SetText(resource.URI)
	itemUri.Disable()

	passwdLabel := widget.NewLabelWithStyle(
		"Password",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemPasswd := widget.NewPasswordEntry()
	itemPasswd.SetPlaceHolder("eg. ************")

	itemPasswd.SetText(resource.Password)
	itemPasswd.Disable()

	descLabel := widget.NewLabelWithStyle(
		"Description",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemDesc := widget.NewEntry()
	itemDesc.SetPlaceHolder("eg. An Amazon account")

	itemDesc.SetText(resource.Description)
	itemDesc.Disable()

	copyUsername := widget.NewButton("Copy Username", func() {
		clipboard.Write(clipboard.FmtText, []byte(resource.Username))
	})

	copyPasswd := widget.NewButton("Copy Password", func() {
		clipboard.Write(clipboard.FmtText, []byte(resource.Password))
	})

	editBtn := widget.NewButton("Edit", func() {
		view.Window.Close()

		editView := NewResourceEditView(ctx, token, resource)
		editView.Window.Show()
	})

	closeBtn := widget.NewButton("Close", func() {
		view.Window.Close()
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
		container.New(
			layout.NewGridLayout(2),
			copyUsername,
			copyPasswd,
		),
		editBtn,
		widget.NewSeparator(),
		closeBtn,
	)
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(350, 100))
	view.Window.CenterOnScreen()
	return view
}
