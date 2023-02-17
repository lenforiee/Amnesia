package main

import (
	"context"

	"fyne.io/fyne/v2/app"

	"github.com/lenforiee/AmnesiaGUI/bundle"
	"github.com/lenforiee/AmnesiaGUI/internals/controllers"
	"github.com/lenforiee/AmnesiaGUI/internals/views"
	"github.com/lenforiee/AmnesiaGUI/models"
	"github.com/lenforiee/AmnesiaGUI/utils"
	"github.com/lenforiee/AmnesiaGUI/utils/logger"
)

var (
	appName = "Amnesia"
)

func main() {
	logger.InitialiseLogging("logs.log")
	logger.LogInfo.Printf("Initialised logger!")

	if err := utils.CheckPidFile(); err != nil {
		logger.LogErr.Printf("Error checking pid file: %s", err)
		return
	}

	if err := utils.NewPidFile(); err != nil {
		logger.LogErr.Printf("Error creating pid file: %s", err)
		return
	}

	logger.LogInfo.Printf("Initialising Amnesia app...")
	app := app.New()
	app.Settings().SetTheme(&models.Theme{})
	app.SetIcon(bundle.ResourceAssetsImagesLogoPng)

	mainWindow := app.NewWindow(appName)
	mainWindow.SetMaster()
	context := context.TODO()

	appContext := controllers.NewAppContext(appName, &app, &mainWindow, &context)
	appContext.LoadConfig()
	appContext.InitialiseSystemTray(bundle.ResourceAssetsImagesLogoPng)

	window, size := views.NewLoginWindow(appContext)
	appContext.UpdateMainWindow(window.Window, size)
	logger.LogInfo.Printf("App initialised!")

	appContext.StartMainWindow()
	utils.RemovePidFile()
}
