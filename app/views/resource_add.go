package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/models"
	"github.com/lenforiee/AmnesiaGUI/app/usecases/passbolt"
)

type ResourceAddView struct {
	Window fyne.Window

	// allow us to have some action in view itself.
	OnButtonBefore func()
	OnButtonError  func()
	OnButtonClick  func()

	Button    *widget.Button
	Container *fyne.Container
}

func NewResourceAddView(ctx amnesiaApp.AppContext) ResourceAddView {

	window := ctx.App.NewWindow(fmt.Sprintf("%s :: Add Resource", ctx.AppName))

	view := ResourceAddView{
		Window: window,
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

		view.OnButtonBefore()

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
			view.OnButtonError()
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
		view.Window.Close()
	})
	view.Button = submitBtn

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
		submitBtn,
		widget.NewSeparator(),
		closeBtn,
	)
	view.Container = containerBox

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(350, 100))
	view.Window.CenterOnScreen()
	return view
}

func (v *ResourceAddView) SetOnButtonBeforeEvent(callback func()) {
	v.OnButtonBefore = callback
}

func (v *ResourceAddView) SetOnButtonErrorEvent(callback func()) {
	v.OnButtonError = callback
}

func (v *ResourceAddView) SetOnButtonClickEvent(callback func()) {
	v.OnButtonClick = callback
}
