package service

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type MessageService struct{}

// 信息对话框
func (g *MessageService) Dialog() {
	app := application.Get()
	if app == nil {
		return
	}
	dialog := app.Dialog.Info()
	dialog.SetTitle("Welcome")
	dialog.SetMessage("Welcome to our application!")
	dialog.Show()
}
