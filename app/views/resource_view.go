package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/models"
)

type ResourceView struct {
	Title     string
	Container *fyne.Container
}

func NewResourceView(ctx *amnesiaApp.AppContext, token string, resource models.Resource, previousView ListView) ResourceView {

	logger.LogInfo.Printf("Creating new resource view for id %s and name %s", token, resource.Name)
	title := fmt.Sprintf("%s :: View Resource", ctx.AppName)

	view := ResourceView{
		Title: title,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name (*)",
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
		"Password (*)",
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
		logger.LogInfo.Printf("Copying resource `%s` username to clipboard", resource.Name)
		clipboard.Write(clipboard.FmtText, []byte(resource.Username))
	})

	copyPasswd := widget.NewButton("Copy Password", func() {
		logger.LogInfo.Printf("Copying resource `%s` password to clipboard", resource.Name)
		clipboard.Write(clipboard.FmtText, []byte(resource.Password))
	})

	editBtn := widget.NewButton("Edit", func() {

		editView := NewResourceEditView(ctx, token, resource, view, previousView)
		ctx.UpdateView(editView.Title, editView.Container)
	})

	goBackBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		ctx.UpdateMainWindow(previousView.Window, previousView.Size, false)
	})

	asteriskLabel := widget.NewLabel("(*) - Required field.")

	containerBox := container.New(
		layout.NewPaddedLayout(),
		container.NewVBox(
			container.NewGridWithColumns(
				6,
				goBackBtn,
			),
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
			asteriskLabel,
			container.New(
				layout.NewGridLayout(2),
				copyUsername,
				copyPasswd,
			),
			editBtn,
		),
	)
	view.Container = containerBox

	return view
}
