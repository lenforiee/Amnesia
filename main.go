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

	amnesiaApp "github.com/lenforiee/Amnesia/app"
	"github.com/lenforiee/Amnesia/app/internals/logger"
	"github.com/lenforiee/Amnesia/app/internals/settings"
	"github.com/lenforiee/Amnesia/app/themes"
	"github.com/lenforiee/Amnesia/app/views"
	"github.com/lenforiee/Amnesia/bundles"
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

	logger.LogInfo.Print("Initialising Amnesia app...")
	app := InitialiseFyneApp()
	logger.LogInfo.Print("Done!")

	defer func() {
		if r := recover(); r != nil {
			for _, win := range app.Driver().AllWindows() {
				logger.LogWarn.Printf("Closing window: %s", win.Title())
				win.Close()
			}
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

	ctx.SetApp(app)

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
	mainView := views.NewLoginView(&ctx)
	ctx.SetMainWindow(mainView.Window)
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("Initialising system tray...")
	ctx.InitialiseSystemTray()
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("Setting close intercept...")
	ctx.MainWindow.SetCloseIntercept(func() {
		if ctx.PassboltClient != nil {
			logger.LogInfo.Print("Hidding main window to tray...")
			ctx.MainWindow.Hide()
		} else {
			logger.LogInfo.Print("Closing main window...")
			ctx.MainWindow.Close()
		}
	})
	logger.LogInfo.Print("Done!")

	logger.LogInfo.Print("App initialised!")
	ctx.MainWindow.ShowAndRun()
}
