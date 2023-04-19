package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/models"
	"github.com/lenforiee/AmnesiaGUI/app/usecases/passbolt"
)

type ResourceEditView struct {
	Title     string
	Container *fyne.Container
}

func NewResourceEditView(
	ctx *amnesiaApp.AppContext,
	token string,
	resource models.Resource,
	previousView ResourceView,
	listView ListView,
) ResourceEditView {

	logger.LogInfo.Printf("Creating new resource edit view for id %s and name %s", token, resource.Name)
	title := fmt.Sprintf("%s :: Edit Resource", ctx.AppName)

	view := ResourceEditView{
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

	saveBtn := widget.NewButton("Save", nil)
	deleteBtn := widget.NewButton("Delete", nil)

	goBackBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		ctx.UpdateView(previousView.Title, previousView.Container)
	})

	deleteBtn.OnTapped = func() {
		saveBtn.Disable()
		deleteBtn.Disable()
		goBackBtn.Disable()

		confirmView := NewConfirmView(ctx, "Are you sure you want to delete this resource?")
		confirmView.SetOnYesEvent(func() {

			err := passbolt.DeleteResource(ctx, token)
			if err != nil {
				errMsg := fmt.Sprintf("There was error while deleting resource: %s", err)
				logger.LogErr.Println(errMsg)

				errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
				errView.Window.Show()
				return
			}

			// views/list.go
			RefreshListData(ctx)
			ctx.UpdateMainWindow(listView.Window, listView.Size, false)
		})

		confirmView.SetOnNoEvent(func() {
			saveBtn.Enable()
			deleteBtn.Enable()
			goBackBtn.Enable()
		})

		confirmView.Window.Show()
	}

	saveBtn.OnTapped = func() {
		saveBtn.Disable()
		deleteBtn.Disable()
		goBackBtn.Disable()
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
		ctx.UpdateMainWindow(listView.Window, listView.Size, false)
	}

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
			saveBtn,
			deleteBtn,
		),
	)

	view.Container = containerBox
	return view
}
