package main

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/lenforiee/PassboltGUI/internals/controllers"
	"github.com/lenforiee/PassboltGUI/internals/views"
	"github.com/lenforiee/PassboltGUI/utils"
	"github.com/lenforiee/PassboltGUI/utils/logger"
)

func main() {
	logger.InitialiseLogging("logs.log")
	logger.LogInfo.Printf("Initialised logger!")

	if err := utils.CheckPidFile(); err != nil {
		logger.LogErr.Printf("Error checking pid file: %s", err)
		return
	}

	logo, err := fyne.LoadResourceFromPath("./assets/logo.png")
	if err != nil {
		logger.LogErr.Printf("Error loading logo: %s", err)
		return
	}

	if err := utils.NewPidFile(); err != nil {
		logger.LogErr.Printf("Error creating pid file: %s", err)
		return
	}

	logger.LogInfo.Printf("Initialising Passbolt app...")
	app := app.New()
	app.SetIcon(logo)

	mainWindow := app.NewWindow("PassboltGUI")
	mainWindow.SetMaster()
	context := context.TODO()

	appContext := controllers.NewAppContext(&app, &mainWindow, &context)
	appContext.LoadConfig()
	appContext.InitialiseSystemTray(logo)

	window, size := views.NewLoginWindow(appContext)
	appContext.UpdateMainWindow(window.Window, size)
	logger.LogInfo.Printf("App initialised!")

	appContext.StartMainWindow()
	utils.RemovePidFile()
}
