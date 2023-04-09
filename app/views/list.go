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
	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/usecases/passbolt"
	"github.com/lenforiee/AmnesiaGUI/bundles"
	"github.com/passbolt/go-passbolt/api"
)

type ListView struct {
	Window    fyne.Window
	Container *fyne.Container
}

var (
	protecetedNameList  = []string{}
	protectedTokenIdMap = make(map[string]interface{})
	lookupNameList      = binding.NewStringList()
	lookupTokenMap      = binding.NewUntypedMap()

	formattedNameList = binding.NewStringList()

	resourcesList []api.Resource

	listWidget   *widget.List
	searchWidget *widget.Entry
)

func NewListView(ctx amnesiaApp.AppContext) ListView {

	window := ctx.App.NewWindow(fmt.Sprintf("%s :: Account List", ctx.AppName))
	view := ListView{
		Window: window,
	}

	resources, err := passbolt.GetResources(ctx, api.GetResourcesOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("There was error while getting user accounts: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
		errView.Window.Show()
	}

	for i, r := range resources {
		protecetedNameList = append(protecetedNameList, r.Name)
		protectedTokenIdMap[strconv.Itoa(i)] = r.ID
		if r.Username == "" {
			formattedNameList.Append(fmt.Sprintf("%s", r.Name))
		} else {
			username := r.Username
			if len(username) > 16 {
				username = fmt.Sprintf("%s...", strings.Split(r.Username, "")[:16])
			}
			formattedNameList.Append(fmt.Sprintf("%s - %s", r.Name, username))
		}
	}

	lookupNameList.Set(protecetedNameList)
	lookupTokenMap.Set(protectedTokenIdMap)
	resourcesList = resources

	list := widget.NewListWithData(formattedNameList,
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

		loadingSplash := NewLoadingSplash(ctx, "Loading resource...")
		loadingSplash.Window.Show()

		token, err := lookupTokenMap.GetValue(strconv.Itoa(id))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting token value: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			loadingSplash.Close()
			return
		}

		resource, err := passbolt.GetResource(ctx, token.(string))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting resource data: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			loadingSplash.Close()
			return
		}

		resourceView := NewResourceView(ctx, token.(string), resource)
		resourceView.Window.Show()
		loadingSplash.Close()

	}
	listWidget = list

	search := widget.NewEntry()
	search.SetPlaceHolder("eg. Amazon")

	search.OnChanged = func(s string) {

		if strings.TrimSpace(s) == "" {
			lookupNameList.Set(protecetedNameList)
			lookupTokenMap.Set(protectedTokenIdMap)
			list.Refresh()
			return
		}

		var filteredData = []string{}
		var filteredTokenIdMap = make(map[string]interface{})
		for _, r := range resourcesList {
			if strings.Contains(strings.ToLower(r.Name), strings.ToLower(s)) {
				filteredData = append(filteredData, r.Name)
				filteredTokenIdMap[strconv.Itoa(len(filteredData)-1)] = r.ID
			}
		}

		lookupNameList.Set(filteredData)
		lookupTokenMap.Set(filteredTokenIdMap)
		list.Refresh()
	}
	searchWidget = search

	hideBtn := widget.NewButton("Hide to tray", func() {
		ctx.MainWindow.Hide()
	})

	// create refresh button
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		loadingSplash := NewLoadingSplash(ctx, "Refreshing the list...")
		loadingSplash.Window.Show()
		RefreshListData(ctx)
		loadingSplash.Close()
	})

	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		addView := NewResourceAddView(ctx)

		loadingSplash := NewLoadingSplash(ctx, "Adding the resource...")
		addView.SetOnButtonBeforeEvent(func() {
			loadingSplash.Window.Show()
		})

		addView.SetOnButtonErrorEvent(func() {
			loadingSplash.Close()
		})

		addView.SetOnButtonClickEvent(func() {
			loadingSplash.UpdateText("Refreshing the list...")
			RefreshListData(ctx)
			loadingSplash.Close()
		})

		addView.Window.Show()
	})

	image := canvas.NewImageFromResource(bundles.ResourceAmnesiaLogoPng)
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

	view.Window.SetContent(view.Container)
	view.Window.Resize(fyne.NewSize(400, 600))
	view.Window.CenterOnScreen()

	return view
}

func RefreshListData(ctx amnesiaApp.AppContext) {
	resources, err := passbolt.GetResources(ctx, api.GetResourcesOptions{})

	if err != nil {
		errMsg := fmt.Sprintf("There was error while refreshing the list: %s", err)
		logger.LogErr.Println(errMsg)

		errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
		errView.Window.Show()
	}

	protecetedNameList = []string{}
	protectedTokenIdMap = make(map[string]interface{})
	resourcesList = resources

	for i, r := range resourcesList {
		protecetedNameList = append(protecetedNameList, r.Name)
		protectedTokenIdMap[strconv.Itoa(i)] = r.ID
	}

	lookupNameList.Set(protecetedNameList)
	lookupTokenMap.Set(protectedTokenIdMap)

	searchWidget.SetText("")
	listWidget.Refresh()

}
