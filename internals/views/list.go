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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lenforiee/AmnesiaGUI/bundle"
	"github.com/lenforiee/AmnesiaGUI/internals/controllers"
	"github.com/lenforiee/AmnesiaGUI/utils/logger"
	"github.com/lenforiee/AmnesiaGUI/utils/passbolt"
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

	List   *widget.List
	Search *widget.Entry
)

func NewListWindow(app *controllers.AppContext) (*ListWindow, fyne.Size) {

	window := (*app.App).NewWindow(fmt.Sprintf("%s :: Account List", app.AppName))
	view := &ListWindow{
		Window:    &window,
		Container: nil,
	}

	resources, err := passbolt.GetResources(app, api.GetResourcesOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("There was error while getting user accounts: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorWindow(app, errMsg)
		app.CreateNewWindowAndShow(errView.Window)
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

		loadingSplash := NewLoadingWindow(app, "Loading resource...")
		app.CreateNewWindowAndShow(loadingSplash.Window)

		token, err := LookupTokenMap.GetValue(strconv.Itoa(id))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting token value: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowAndShow(errView.Window)
			loadingSplash.StopLoading(app)
			return
		}

		resource, err := passbolt.GetResource(app, token.(string))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting resource data: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorWindow(app, errMsg)
			app.CreateNewWindowAndShow(errView.Window)
			loadingSplash.StopLoading(app)
			return
		}

		resourceView := NewResourceWindow(app, token.(string), resource)
		app.CreateNewWindowAndShow(resourceView.Window)
		loadingSplash.StopLoading(app)

	}
	List = list

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
	Search = search

	hideBtn := widget.NewButton("Hide to tray", func() {
		(*app.MainWindow).Hide()
	})

	// create refresh button
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		loadingSplash := NewLoadingWindow(app, "Refreshing the list...")
		app.CreateNewWindowAndShow(loadingSplash.Window)
		RefreshListData(app)
		loadingSplash.StopLoading(app)
	})

	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		addView := NewResourceAddWindow(app)

		loadingSplash := NewLoadingWindow(app, "Adding the resource...")
		addView.OnButtonBefore = func() {
			app.CreateNewWindowAndShow(loadingSplash.Window)
		}

		addView.OnButtonError = func() {
			loadingSplash.StopLoading(app)
		}

		addView.OnButtonClick = func() {
			loadingSplash.UpdateText("Refreshing the list...")
			RefreshListData(app)
			loadingSplash.StopLoading(app)
		}

		app.CreateNewWindowAndShow(addView.Window)
	})

	image := canvas.NewImageFromResource(bundle.ResourceAssetsImagesAmnesialogoPng)
	image.FillMode = canvas.ImageFillOriginal

	containerBox := container.NewBorder(
		container.New(
			layout.NewVBoxLayout(),
			image,
			container.NewBorder(
				nil,
				nil,
				nil,
				container.New(
					layout.NewGridLayout(2),
					addBtn,
					refreshBtn,
				),
				search,
			),
		),
		container.New(
			layout.NewVBoxLayout(),
			hideBtn,
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

func RefreshListData(app *controllers.AppContext) {
	resources, err := passbolt.GetResources(app, api.GetResourcesOptions{})

	if err != nil {
		errMsg := fmt.Sprintf("There was error while refreshing the list: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorWindow(app, errMsg)
		app.CreateNewWindowAndShow(errView.Window)
	}

	RealNameList = []string{}
	RealTokenIdMap = make(map[string]interface{})
	Resources = &resources

	for i, r := range *Resources {
		RealNameList = append(RealNameList, r.Name)
		RealTokenIdMap[strconv.Itoa(i)] = r.ID
	}

	LookupNameList.Set(RealNameList)
	LookupTokenMap.Set(RealTokenIdMap)

	Search.SetText("")
	List.Refresh()

}
