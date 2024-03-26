package controller

import (
	"github.com/gin-gonic/gin"
	wapi "github.com/iamacarpet/go-win64api"
	"github.com/prometheus-community/windows_exporter/custom/bean/vo/request"
	"github.com/prometheus-community/windows_exporter/custom/resp"
)

func SetAdmin(c *gin.Context) {
	var req request.SetAdmin
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.SetAdmin(req.UserName)
	commonReturn(c, err, ok)
}

func RevokeAdmin(c *gin.Context) {
	var req request.SetAdmin
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.RevokeAdmin(req.UserName)
	commonReturn(c, err, ok)
}
