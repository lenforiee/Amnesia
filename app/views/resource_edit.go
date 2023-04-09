package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/models"
	"github.com/lenforiee/AmnesiaGUI/app/usecases/passbolt"
)

type ResourceEditView struct {
	Window    fyne.Window
	Container *fyne.Container
}

// TODO: fix this function, it consistantly crashes the app.
func NewResourceEditView(ctx app.AppContext, token string, resource models.Resource) ResourceEditView {

	window := ctx.App.NewWindow(fmt.Sprintf("%s :: Edit Resource", ctx.AppName))

	view := ResourceEditView{
		Window: window,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name (*)",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemName := widget.NewEntry()
	itemName.SetPlaceHolder("eg. Amazon")

	itemName.SetText(resource.Name)

	usernameLabel := widget.NewLabelWithStyle(
		"Username",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemUsername := widget.NewEntry()
	itemUsername.SetPlaceHolder("eg. example@example.com")

	itemUsername.SetText(resource.Username)

	uriLabel := widget.NewLabelWithStyle(
		"URI",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemUri := widget.NewEntry()
	itemUri.SetPlaceHolder("eg. https://amazon.com")

	itemUri.SetText(resource.URI)

	passwdLabel := widget.NewLabelWithStyle(
		"Password (*)",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemPasswd := widget.NewPasswordEntry()
	itemPasswd.SetPlaceHolder("eg. ************")

	itemPasswd.SetText(resource.Password)

	descLabel := widget.NewLabelWithStyle(
		"Description",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemDesc := widget.NewEntry()
	itemDesc.SetPlaceHolder("eg. An Amazon account")

	itemDesc.SetText(resource.Description)

	asteriskLabel := widget.NewLabel("(*) - Required field.")

	saveBtn := widget.NewButton("Save", func() {
		saveResource := models.NewResource()
		saveResource.SetFolderParentID(resource.FolderParentID)
		saveResource.SetName(itemName.Text)
		saveResource.SetUsername(itemUsername.Text)
		saveResource.SetURI(itemUri.Text)
		saveResource.SetPassword(itemPasswd.Text)
		saveResource.SetDescription(itemDesc.Text)

		err := passbolt.UpdateResource(ctx, token, saveResource)
		if err != nil {
			errMsg := fmt.Sprintf("There was error while updating resource: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			return
		}

		// views/list.go
		RefreshListData(ctx)
		view.Window.Close()
	})

	deleteBtn := widget.NewButton("Delete", func() {
		confirmView := NewConfirmView(ctx, "Are you sure you want to delete this resource?")
		confirmView.SetOnYesEvent(func() {
			loadingSplash := NewLoadingSplash(ctx, "Removing the resource...")
			loadingSplash.Window.Show()

			err := passbolt.DeleteResource(ctx, token)
			if err != nil {
				errMsg := fmt.Sprintf("There was error while deleting resource: %s", err)
				logger.LogErr.Println(errMsg)

				errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
				errView.Window.Show()
				loadingSplash.Close()
				return
			}

			// views/list.go
			loadingSplash.UpdateText("Refreshing the list...")
			RefreshListData(ctx)
			loadingSplash.Close()
		})

		confirmView.Window.Show()
		view.Window.Close()
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
		asteriskLabel,
		saveBtn,
		deleteBtn,
		widget.NewSeparator(),
		closeBtn,
	)
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(350, 100))
	view.Window.CenterOnScreen()
	return view
}
