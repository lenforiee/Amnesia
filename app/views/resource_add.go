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

type ResourceAddView struct {
	Title string

	// allow us to have some action in view itself.
	OnButtonClick func()

	Button    *widget.Button
	Container *fyne.Container
}

func NewResourceAddView(ctx *amnesiaApp.AppContext, previousView ListView) *ResourceAddView {

	logger.LogInfo.Println("Creating new resource add view")
	title := fmt.Sprintf("%s :: Add Resource", ctx.AppName)

	view := &ResourceAddView{
		Title: title,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name (*)",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	itemName := widget.NewEntry()
	itemName.SetPlaceHolder("eg. Amazon")

	usernameLabel := widget.NewLabelWithStyle(
		"Username",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemUsername := widget.NewEntry()
	itemUsername.SetPlaceHolder("eg. example@example.com")

	uriLabel := widget.NewLabelWithStyle(
		"URI",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemUri := widget.NewEntry()
	itemUri.SetPlaceHolder("eg. https://amazon.com")

	passwdLabel := widget.NewLabelWithStyle(
		"Password (*)",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemPasswd := widget.NewPasswordEntry()
	itemPasswd.SetPlaceHolder("eg. ************")

	descLabel := widget.NewLabelWithStyle(
		"Description",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	itemDesc := widget.NewEntry()
	itemDesc.SetPlaceHolder("eg. An Amazon account")

	asteriskLabel := widget.NewLabel("(*) - Required field.")

	submitBtn := widget.NewButton("Submit", func() {

		var emptyFields []string

		if itemName.Text == "" {
			emptyFields = append(emptyFields, "Resource Name")
		}

		if itemPasswd.Text == "" {
			emptyFields = append(emptyFields, "Password")
		}

		if len(emptyFields) > 0 {

			bulletPoints := ""
			for _, field := range emptyFields {
				bulletPoints += fmt.Sprintf(" • %s\n", field)
			}

			errMsg := fmt.Sprintf("The following fields are empty: \n%sPlease fill them to continue.", bulletPoints)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			return
		}

		err := passbolt.CreateResource(ctx,
			models.Resource{
				FolderParentID: "",
				Name:           itemName.Text,
				Username:       itemUsername.Text,
				URI:            itemUri.Text,
				Password:       itemPasswd.Text,
				Description:    itemDesc.Text,
			},
		)

		if err != nil {
			errMsg := fmt.Sprintf("There was error while adding a resource: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			return
		}

		view.OnButtonClick()
		ctx.UpdateMainWindow(previousView.Window, previousView.Size, false)
	})
	view.Button = submitBtn

	goBackBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		ctx.UpdateMainWindow(previousView.Window, previousView.Size, false)
	})

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
			submitBtn,
		),
	)

	view.Container = containerBox
	return view
}

func (v *ResourceAddView) SetOnButtonClickEvent(callback func()) {
	v.OnButtonClick = callback
}
