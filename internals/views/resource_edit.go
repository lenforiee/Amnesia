package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/lenforiee/AmnesiaGUI/internals/controllers"
	"github.com/lenforiee/AmnesiaGUI/models"
	"github.com/lenforiee/AmnesiaGUI/utils/logger"
	"github.com/lenforiee/AmnesiaGUI/utils/passbolt"
)

type ResourceEditWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

// TODO: fix this function, it consistantly crashes the app.
func NewResourceEditWindow(app *controllers.AppContext, token string, resource *models.Resource) *ResourceEditWindow {

	window := (*app.App).NewWindow(fmt.Sprintf("%s :: Edit Resource", app.AppName))

	view := &ResourceEditWindow{
		Window:    &window,
		Container: nil,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name",
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
		"Password",
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

	saveBtn := widget.NewButton("Save", func() {
		saveResource := models.Resource{
			FolderParentID: "",
			Name:           itemName.Text,
			Username:       itemUsername.Text,
			URI:            itemUri.Text,
			Password:       itemPasswd.Text,
			Description:    itemDesc.Text,
		}
		err := passbolt.UpdateResource(app, token, &saveResource)
		if err != nil {
			errMsg := fmt.Sprintf("There was error while updating resource: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowAndShow(errView.Window)
			return
		}

		// views/list.go
		RefreshListData(app)

		(*view.Window).Close()
	})

	deleteBtn := widget.NewButton("Delete", func() {
		confirmView := NewConfirmWindow(app, "Are you sure you want to delete this resource?")
		confirmView.OnYes = func() {
			err := passbolt.DeleteResource(app, token)
			if err != nil {
				errMsg := fmt.Sprintf("There was error while deleting resource: %s", err)
				logger.LogErr.Println(errMsg)

				errView := NewErrorWindow(app, errMsg)
				app.CreateNewWindowAndShow(errView.Window)
				return
			}

			// views/list.go
			loadingSplash := NewLoadingWindow(app, "Refreshing the list...")
			app.CreateNewWindowAndShow(loadingSplash.Window)
			RefreshListData(app)
			loadingSplash.StopLoading(app)
		}

		app.CreateNewWindowAndShow(confirmView.Window)
		(*view.Window).Close()
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
		saveBtn,
		deleteBtn,
		widget.NewSeparator(),
		closeBtn,
	)
	view.Container = containerBox

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(fyne.NewSize(350, 100))
	(*view.Window).CenterOnScreen()
	return view
}
