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
	"github.com/lenforiee/AmnesiaGUI/app/usecases/format"
	"github.com/lenforiee/AmnesiaGUI/app/usecases/passbolt"
	"github.com/lenforiee/AmnesiaGUI/bundles"
	"github.com/passbolt/go-passbolt/api"
)

type ListView struct {
	Window fyne.Window

	Size      fyne.Size
	Container *fyne.Container
}

var (
	protecetedNameList         = []string{}
	protectedTokenIdMap        = make(map[string]interface{})
	protectedFormattedNameList = []string{}
	lookupNameList             = binding.NewStringList()
	lookupTokenMap             = binding.NewUntypedMap()
	formattedNameList          = binding.NewStringList()

	resourcesList []api.Resource

	refreshBtnWidget *widget.Button
	listWidget       *widget.List
	searchWidget     *widget.Entry
)

func NewListView(ctx *amnesiaApp.AppContext) ListView {

	logger.LogInfo.Println("Creating new list view")
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
			protectedFormattedNameList = append(protectedFormattedNameList, r.Name)
		} else {
			protectedFormattedNameList = append(
				protectedFormattedNameList,
				fmt.Sprintf("%s - (%s)", r.Name, format.TruncateText(r.Username, 32)),
			)
		}
	}

	lookupNameList.Set(protecetedNameList)
	lookupTokenMap.Set(protectedTokenIdMap)
	formattedNameList.Set(protectedFormattedNameList)
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

		token, err := lookupTokenMap.GetValue(strconv.Itoa(id))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting token value: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			return
		}

		resource, err := passbolt.GetResource(ctx, token.(string))
		if err != nil {
			errMsg := fmt.Sprintf("There was error while getting resource data: %s", err)
			logger.LogErr.Println(errMsg)

			errView := NewErrorView(ctx.App, ctx.AppName, errMsg, false)
			errView.Window.Show()
			return
		}

		resourceView := NewResourceView(ctx, token.(string), resource, view)
		ctx.UpdateView(resourceView.Title, resourceView.Container)

	}
	listWidget = list

	search := widget.NewEntry()
	search.SetPlaceHolder("eg. Amazon")

	search.OnChanged = func(s string) {

		logger.LogInfo.Printf("Search query: '%s'", s)
		if strings.TrimSpace(s) == "" {
			logger.LogInfo.Print("Search query is empty, resetting list")
			lookupNameList.Set(protecetedNameList)
			lookupTokenMap.Set(protectedTokenIdMap)
			formattedNameList.Set(protectedFormattedNameList)
			list.Refresh()
			return
		}

		var filteredData = []string{}
		var filteredTokenIdMap = make(map[string]interface{})
		var filteredFormattedNameList = []string{}
		for _, r := range resourcesList {
			if strings.Contains(strings.ToLower(r.Name), strings.ToLower(s)) {
				filteredData = append(filteredData, r.Name)
				filteredTokenIdMap[strconv.Itoa(len(filteredData)-1)] = r.ID

				if r.Username == "" {
					filteredFormattedNameList = append(filteredFormattedNameList, r.Name)
				} else {
					filteredFormattedNameList = append(
						filteredFormattedNameList,
						fmt.Sprintf("%s - (%s)", r.Name, format.TruncateText(r.Username, 32)),
					)
				}
			}
		}

		lookupNameList.Set(filteredData)
		lookupTokenMap.Set(filteredTokenIdMap)
		formattedNameList.Set(filteredFormattedNameList)

		list.Refresh()
	}
	searchWidget = search

	// create refresh button
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		RefreshListData(ctx)
	})
	refreshBtnWidget = refreshBtn

	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		addView := NewResourceAddView(ctx, view)

		addView.SetOnButtonClickEvent(func() {
			RefreshListData(ctx)
		})

		ctx.UpdateView(addView.Title, addView.Container)
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
		nil,
		nil,
		nil,
		list,
	)
	view.Container = containerBox
	view.Size = fyne.NewSize(400, 550)

	view.Window.SetContent(view.Container)
	view.Window.Resize(view.Size)
	view.Window.CenterOnScreen()

	return view
}

func RefreshListData(ctx *amnesiaApp.AppContext) {

	logger.LogInfo.Println("Refreshing list data")
	refreshBtnWidget.Disable()
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
	newFormattedNameList := []string{}

	for i, r := range resourcesList {
		protecetedNameList = append(protecetedNameList, r.Name)
		protectedTokenIdMap[strconv.Itoa(i)] = r.ID

		if r.Username == "" {
			newFormattedNameList = append(newFormattedNameList, r.Name)
		} else {
			newFormattedNameList = append(
				newFormattedNameList,
				fmt.Sprintf("%s - (%s)", r.Name, format.TruncateText(r.Username, 32)),
			)
		}
	}

	lookupNameList.Set(protecetedNameList)
	lookupTokenMap.Set(protectedTokenIdMap)
	protectedFormattedNameList = newFormattedNameList
	formattedNameList.Set(protectedFormattedNameList)

	searchWidget.SetText("")
	listWidget.Refresh()
	refreshBtnWidget.Enable()
}
