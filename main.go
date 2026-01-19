package main

import (
	"context"
	"embed"
	"os"
	stdruntime "runtime"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed resources/icon.ico
var iconData []byte

//go:embed resources/icon.png
var pngIconData []byte

var app *App

func main() {
	// Create an instance of the app structure
	app = NewApp()
	app.BaseIcon = pngIconData

	// Start systray in a goroutine
	go func() {
		systray.Run(onReady, onExit)
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Sungrow iSolarCloud Monitor",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		OnBeforeClose: func(ctx context.Context) bool {
			// Hide window instead of closing
			runtime.WindowHide(ctx)
			return true // Prevent close
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func onReady() {
	// Set icon - fyne.io/systray has better Windows support
	if stdruntime.GOOS == "windows" {
		systray.SetIcon(iconData)
	} else {
		systray.SetIcon(pngIconData)
	}
	systray.SetTitle("Sungrow")
	systray.SetTooltip("Sungrow iSolarCloud Monitor")

	// Add menu items
	mShow := systray.AddMenuItem("Show App", "Show the main window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Listen for tray title updates
	go func() {
		for title := range app.TrayTitleChan {
			systray.SetTooltip(title)
		}
	}()

	// Listen for tray icon updates
	go func() {
		for iconBytes := range app.TrayIconChan {
			println("Received tray icon update, length:", len(iconBytes))
			systray.SetIcon(iconBytes)
		}
	}()

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if app.ctx != nil {
					runtime.WindowShow(app.ctx)
					runtime.WindowUnminimise(app.ctx)
				}
			case <-mQuit.ClickedCh:
				// Properly quit everything
				systray.Quit()
				if app.ctx != nil {
					runtime.Quit(app.ctx)
				}
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	// Cleanup
	os.Exit(0)
}
