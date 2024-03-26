package controller

import (
	"github.com/gin-gonic/gin"
	wapi "github.com/iamacarpet/go-win64api"
	so "github.com/iamacarpet/go-win64api/shared"
	"github.com/prometheus-community/windows_exporter/custom/bean/vo/request"
	"github.com/prometheus-community/windows_exporter/custom/resp"
)

func AddWindowsLocalGroup(c *gin.Context) {
	var req request.AddWindowsLocalGroup
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.LocalGroupAdd(req.GroupName, req.Comment)
	commonReturn(c, err, ok)
}

func DeleteWindowsLocalGroup(c *gin.Context) {
	var req request.DeleteWindowsLocalGroup
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.LocalGroupDel(req.GroupName)
	commonReturn(c, err, ok)
}

func ListWindowsLocalGroup(c *gin.Context) {
	var req request.ListWindowsLocalGroup
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	data, err := wapi.ListLocalGroups()
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, data)
}

func DetailWindowsLocalGroup(c *gin.Context) {
	var req request.DetailWindowsLocalGroup
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	localGroups, err := wapi.ListLocalGroups()
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}

	var data *so.LocalGroup
	for _, v := range localGroups {
		if v.Name == req.GroupName {
			data = &v
			break
		}
	}

	if data == nil {
		resp.SendJsonResponse(c, resp.NotFound.WithMsg("数据不存在"), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, data)
}
