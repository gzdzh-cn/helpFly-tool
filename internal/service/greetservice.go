package service

import (
	gf "github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct{}

var mainWindow *application.WebviewWindow

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}

// 设置主题
func (g *GreetService) SetTheme() {
	gf.Log().Info(gctx.New(), "执行设置主题函数", application.Get().GetPID())
	if mainWindow != nil {
		mainWindow.Restore()
		mainWindow.Focus()
	}
}
