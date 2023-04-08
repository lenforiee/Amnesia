package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	amnesiaApp "github.com/lenforiee/AmnesiaGUI/app"
	"github.com/lenforiee/AmnesiaGUI/app/internals/logger"
	"github.com/lenforiee/AmnesiaGUI/app/internals/settings"
	"github.com/lenforiee/AmnesiaGUI/app/themes"
	"github.com/lenforiee/AmnesiaGUI/app/views"
	"github.com/lenforiee/AmnesiaGUI/bundles"
)

var (
	appName = "Amnesia"
)

func InitialiseFyneApp() fyne.App {
	app := app.NewWithID("com.lenforiee.amnesia")
	app.Settings().SetTheme(themes.ClassicThemeDark())
	app.SetIcon(bundles.ResourceLogoPng)

	return app
}

func main() {

	logger.InitialiseLogging("logs.log")
	logger.LogInfo.Print("Initialised logger!")

	// TODO: check if this code even works.
	defer func() {
		if r := recover(); r != nil {
			app := InitialiseFyneApp()
			logger.LogErr.Printf("Program has crashed: %s", string(debug.Stack()))
			errorView := views.NewErrorView(
				app, appName,
				fmt.Sprintf("Please, report this issue to the developer.\nError: %s", r),
				true, // Crash.
			)
			errorView.Window.ShowAndRun()
		}
	}()

	ctx := amnesiaApp.NewAppContext()
	ctx.SetAppName(appName)

	// checks if port is already occupied (one-instance check)
	logger.LogInfo.Print("Checking if one instance is already running...")
	conn, _ := net.DialTimeout("tcp", "127.0.0.1:44557", time.Second*1)
	if conn != nil {
		app := InitialiseFyneApp()
		mainView := views.NewErrorView(
			app, appName,
			"Amnesia is already running!\nPlease, close the previous instance of Amnesia and try again.",
			false, // Not a crash.
		)
		logger.LogErr.Print("Fail: Amnesia is already running!")
		mainView.Window.ShowAndRun()
		os.Exit(1)
	}
	logger.LogInfo.Print("Done!")

	// initialises a listener on port 44557
	logger.LogInfo.Print("Initialising a port listener...")
	listener, err := net.Listen("tcp", "127.0.0.1:44557")
	if err != nil {
		panic(err.Error()) // Error handler will catch it.
	}
	logger.LogInfo.Print("Done!")
	defer listener.Close()

	logger.LogInfo.Print("Initialising Amnesia app...")
	app := InitialiseFyneApp()
	ctx.SetApp(app)
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("Initialising connector context...")
	context := context.TODO()
	ctx.SetContext(context)
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("Loading user settings...")
	cfg, err := settings.LoadUserSettings()
	if err != nil {
		panic(err.Error())
	}
	ctx.SetUserSettings(cfg)
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("Initialising main view...")
	mainView := views.NewLoginView(ctx)
	ctx.SetMainWindow(mainView.Window)
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("App initialised!")
	ctx.MainWindow.ShowAndRun()
}
