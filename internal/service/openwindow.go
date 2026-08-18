package service

import "github.com/wailsapp/wails/v3/pkg/application"

type OpenWindow struct{}

func (g *OpenWindow) Open(url string) string {
	app := application.Get()
	if app == nil {
		return "应用尚未初始化"
	}

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "新打开的演示窗口",
		Width:  1300,
		Height: 800,
	})
	window.SetURL(url)

	return "窗口打开成功"
}
