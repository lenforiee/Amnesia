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
	amnesiaApp "github.com/lenforiee/Amnesia/app"
	"github.com/lenforiee/Amnesia/app/internals/logger"
	"github.com/lenforiee/Amnesia/app/usecases/format"
	"github.com/lenforiee/Amnesia/app/usecases/passbolt"
	"github.com/lenforiee/Amnesia/bundles"
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

	noResourcesFoundText = "No resources found"
	noResources          = false
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

	noResources = len(resources) == 0
	for i, r := range resources {
		protecetedNameList = append(protecetedNameList, r.Name)
		protectedTokenIdMap[strconv.Itoa(i)] = r.ID
		if r.Username == "" {
			protectedFormattedNameList = append(protectedFormattedNameList, format.TruncateText(r.Name, 45, false))
		} else {
			protectedFormattedNameList = append(
				protectedFormattedNameList,
				format.TruncateText(fmt.Sprintf("%s - (%s)", r.Name, r.Username), 45, true),
			)
		}
	}

	if noResources && len(protectedFormattedNameList) == 0 {
		protectedFormattedNameList = append(protectedFormattedNameList, noResourcesFoundText)
	}

	lookupNameList.Set(protecetedNameList)
	lookupTokenMap.Set(protectedTokenIdMap)
	formattedNameList.Set(protectedFormattedNameList)
	resourcesList = resources

	list := widget.NewListWithData(formattedNameList,
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			label.TextStyle = fyne.TextStyle{Bold: true}
			label.Alignment = fyne.TextAlignCenter

			return label
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		list.Unselect(id)

		if noResources {
			return
		}

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
					filteredFormattedNameList = append(filteredFormattedNameList, format.TruncateText(r.Name, 45, false))
				} else {
					filteredFormattedNameList = append(
						filteredFormattedNameList,
						format.TruncateText(fmt.Sprintf("%s - (%s)", r.Name, r.Username), 45, true),
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

	if noResources {
		search.Disable()
	}

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

	noResources = len(resourcesList) == 0
	for i, r := range resourcesList {
		protecetedNameList = append(protecetedNameList, r.Name)
		protectedTokenIdMap[strconv.Itoa(i)] = r.ID

		if r.Username == "" {
			newFormattedNameList = append(newFormattedNameList, format.TruncateText(r.Name, 45, false))
		} else {
			newFormattedNameList = append(
				newFormattedNameList,
				format.TruncateText(fmt.Sprintf("%s - (%s)", r.Name, r.Username), 45, true),
			)
		}
	}

	if noResources && len(newFormattedNameList) == 0 {
		newFormattedNameList = append(newFormattedNameList, noResourcesFoundText)
	}

	if !noResources {
		searchWidget.Enable()
	}

	lookupNameList.Set(protecetedNameList)
	lookupTokenMap.Set(protectedTokenIdMap)
	protectedFormattedNameList = newFormattedNameList
	formattedNameList.Set(protectedFormattedNameList)

	searchWidget.SetText("")
	listWidget.Refresh()
	refreshBtnWidget.Enable()
}
