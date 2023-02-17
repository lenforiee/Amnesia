package views

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/PassboltGUI/internals/controllers"
	"github.com/lenforiee/PassboltGUI/utils/logger"
	"github.com/lenforiee/PassboltGUI/utils/passbolt"
	"github.com/passbolt/go-passbolt/api"
)

type ListWindow struct {
	Window    *fyne.Window
	Container *fyne.Container
}

var (
	RealNameList   = []string{}
	RealTokenIdMap = make(map[string]interface{})
	LookupNameList = binding.NewStringList()
	LookupTokenMap = binding.NewUntypedMap()
	Resources      *[]api.Resource
)

func NewListWindow(app *controllers.AppContext) (*ListWindow, fyne.Size) {

	window := (*app.App).NewWindow("PassboltGUI Account List")
	view := &ListWindow{
		Window:    &window,
		Container: nil,
	}

	resources, err := passbolt.GetResources(app, api.GetResourcesOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("There was error while getting user accounts: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorWindow(app, errMsg)
		app.CreateNewWindowWithView(errView.Window)
	}

	for i, r := range resources {
		RealNameList = append(RealNameList, r.Name)
		RealTokenIdMap[strconv.Itoa(i)] = r.ID
	}

	LookupNameList.Set(RealNameList)
	LookupTokenMap.Set(RealTokenIdMap)
	Resources = &resources

	list := widget.NewListWithData(LookupNameList,
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			label.TextStyle = fyne.TextStyle{Bold: true}

			return label
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	list.OnSelected = func(id widget.ListItemID) {
		list.Unselect(id)

		token, err := LookupTokenMap.GetValue(strconv.Itoa(id))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting token value: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowWithView(errView.Window)
			return
		}

		resource, err := passbolt.GetResource(app, token.(string))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting resource data: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowWithView(errView.Window)
			return
		}

		resourceView := NewResourceWindow(app, resource)
		app.CreateNewWindowWithView(resourceView.Window)

	}

	search := widget.NewEntry()
	search.SetPlaceHolder("eg. Amazon")

	search.OnChanged = func(s string) {

		if strings.TrimSpace(s) == "" {
			LookupNameList.Set(RealNameList)
			LookupTokenMap.Set(RealTokenIdMap)
			list.Refresh()
			return
		}

		var filteredData = []string{}
		var filteredTokenIdMap = make(map[string]interface{})
		for _, d := range *Resources {
			if strings.Contains(strings.ToLower(d.Name), strings.ToLower(s)) {
				filteredData = append(filteredData, d.Name)
				filteredTokenIdMap[strconv.Itoa(len(filteredData)-1)] = d.ID
			}
		}

		LookupNameList.Set(filteredData)
		LookupTokenMap.Set(filteredTokenIdMap)
		list.Refresh()
	}

	quitBtn := widget.NewButton("Quit", func() {
		(*app.MainWindow).Close()
	})

	// create refresh button
	refreshBtn := widget.NewButton("Refresh List", func() {
		RefreshListData(list, search, app)
	})

	image := canvas.NewImageFromFile("./assets/logo_white.png")
	image.FillMode = canvas.ImageFillOriginal

	containerBox := container.NewBorder(
		image,
		container.New(
			layout.NewVBoxLayout(),
			search,
			refreshBtn,
			widget.NewSeparator(),
			quitBtn,
		),
		nil,
		nil,
		list,
	)
	view.Container = containerBox

	size := fyne.NewSize(400, 600)
	(*view.Window).SetContent(view.Container)
	(*view.Window).Resize(size)
	(*view.Window).CenterOnScreen()

	return view, size
}

func RefreshListData(list *widget.List, search *widget.Entry, app *controllers.AppContext) {
	resources, err := passbolt.GetResources(app, api.GetResourcesOptions{})

	if err != nil {
		errMsg := fmt.Sprintf("There was error while refreshing the list: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorWindow(app, errMsg)
		app.CreateNewWindowWithView(errView.Window)
	}

	RealNameList = []string{}
	RealTokenIdMap = make(map[string]interface{})

	for i, r := range resources {
		RealNameList = append(RealNameList, r.Name)
		RealTokenIdMap[strconv.Itoa(i)] = r.ID
	}

	LookupNameList.Set(RealNameList)
	LookupTokenMap.Set(RealTokenIdMap)

	search.SetText("")
	list.Refresh()

}
