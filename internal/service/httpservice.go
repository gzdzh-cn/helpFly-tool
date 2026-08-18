package service

import (
	"helpfly/internal/service/utils"
)

type HttpService struct{}

// 测试在Service中http-POST请求
func (g *HttpService) PostData(param string) any {
	res, err := utils.PostApi("http://localhost:8200/wails/httprequest/postData", map[string]interface{}{"data": param}, "")
	if err != nil {
		return err.Error()
	}
	return res
}

// 测试在Service中http-GET请求
func (g *HttpService) GetData(param string) any {
	res, err := utils.GetApi("http://localhost:8200/wails/httprequest/getData", map[string]interface{}{"data": param}, "")
	if err != nil {
		return err.Error()
	}
	return res
}
