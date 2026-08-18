package main

import (
	"embed"
	_ "embed"
	"helpfly/internal/service"
	"helpfly/internal/service/db"
	"log"
	"time"

	gf "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	if err := db.Open(gctx.New()); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(gctx.New()); err != nil {
			log.Printf("关闭 SQLite 失败: %v", err)
		}
	}()

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "helpfly",
		Description: "helpfly是一个gf脚手架，使用goframe+wails3 GUI基础框架帮助开发者快速开发桌面应用",
		Services: []application.Service{
			application.NewService(&service.GreetService{}),
			application.NewService(&service.MessageService{}),
			application.NewService(&service.OpenWindow{}),
			application.NewService(&service.HttpService{}),
			application.NewService(&service.DbService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	mainWind := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "helpFly助手(1.0.0)",
		Width:     1280, // 建议初始大小不要小于最小值
		Height:    800,
		MinWidth:  1100, // 最小宽度
		MinHeight: 720,  // 最小高度
		// Frameless: true,//移除窗口的窗框(即边框、头部标题和操作按钮)，用于自定义窗口头部
		BackgroundColour: application.NewRGB(255, 255, 255),
		Windows:          application.WindowsWindow{Theme: 0}, //这里是设框架主题，0=跟随系统，1=Dark(黑色)，2=Light(浅色)
		// macOS 使用原生标题栏，与 Windows 的白色系统顶部栏保持一致。
		// 不使用隐藏/透明标题栏，避免 macOS 顶部栏消失。
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropNormal,
			TitleBar: application.MacTitleBarDefault,
		},
		URL: "/",
	})

	// 监听窗口尺寸变化并通知前端
	mainWind.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		width, height := mainWind.Size()
		app.Event.Emit("window-resize", map[string]int{
			"width":  width,
			"height": height,
		})
	})
	// 首次同步窗口尺寸
	if width, height := mainWind.Size(); width > 0 && height > 0 {
		app.Event.Emit("window-resize", map[string]int{
			"width":  width,
			"height": height,
		})
	}
	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		gf.Log().Info(gctx.New(), "窗口id：", mainWind.ID())
		for {
			now := time.Now().Format(time.DateTime)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()
	//处理监听前端自定义事件
	app.Event.On("myevent", func(e *application.CustomEvent) {
		// Access event information
		gf.Log().Info(gctx.New(), "main监听事件：", e)
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}

}
