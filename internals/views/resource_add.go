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

type ResourceAddWindow struct {
	Window        *fyne.Window
	Button        *widget.Button
	OnButtonClick func()
	Container     *fyne.Container
}

func NewResourceAddWindow(app *controllers.AppContext) *ResourceAddWindow {

	window := (*app.App).NewWindow(fmt.Sprintf("%s :: Add Resource", app.AppName))

	view := &ResourceAddWindow{
		Window:    &window,
		Container: nil,
	}

	nameLabel := widget.NewLabelWithStyle(
		"Resource Name",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	itemName := widget.NewEntry()
	itemName.SetPlaceHolder("eg. Amazon")

	usernameLabel := widget.NewLabelWithStyle(
		"Username",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemUsername := widget.NewEntry()
	itemUsername.SetPlaceHolder("eg. example@example.com")

	uriLabel := widget.NewLabelWithStyle(
		"URI",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemUri := widget.NewEntry()
	itemUri.SetPlaceHolder("eg. https://amazon.com")

	passwdLabel := widget.NewLabelWithStyle(
		"Password",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemPasswd := widget.NewPasswordEntry()
	itemPasswd.SetPlaceHolder("eg. ************")

	descLabel := widget.NewLabelWithStyle(
		"Description",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	itemDesc := widget.NewEntry()
	itemDesc.SetPlaceHolder("eg. An Amazon account")

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

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowAndShow(errView.Window)
			return
		}

		err := passbolt.CreateResource(app,
			&models.Resource{
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

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowAndShow(errView.Window)
			return
		}
		(*view).OnButtonClick()

		(*view.Window).Close()
	})
	view.Button = submitBtn

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
		submitBtn,
		widget.NewSeparator(),
		closeBtn,
	)
	view.Container = containerBox

	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(fyne.NewSize(350, 100))
	(*view.Window).CenterOnScreen()
	return view
}
