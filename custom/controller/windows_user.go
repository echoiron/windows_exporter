package controller

import (
	"github.com/gin-gonic/gin"
	wapi "github.com/iamacarpet/go-win64api"
	so "github.com/iamacarpet/go-win64api/shared"
	"github.com/prometheus-community/windows_exporter/custom/bean/vo/request"
	"github.com/prometheus-community/windows_exporter/custom/resp"
)

func commonReturn(c *gin.Context, err error, ok bool) {
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	if !ok {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg("添加用户失败"), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, nil)
}

func AddWindowsUser(c *gin.Context) {
	var req request.AddWindowsUser
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.UserAdd(req.UserName, req.FullName, req.Password)
	commonReturn(c, err, ok)
}

func UpdateFullNameWindowsUser(c *gin.Context) {
	var req request.UpdateFullNameWindowsUser
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.UserUpdateFullname(req.UserName, req.FullName)
	commonReturn(c, err, ok)
}

func UpdatePasswordWindowsUser(c *gin.Context) {
	var req request.UpdatePasswordWindowsUser
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.ChangePassword(req.UserName, req.Password)
	commonReturn(c, err, ok)
}

func DeleteWindowsUser(c *gin.Context) {
	var req request.DeleteWindowsUser
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	ok, err := wapi.UserDelete(req.UserName)
	commonReturn(c, err, ok)
}

func ListWindowsUser(c *gin.Context) {
	data, err := wapi.ListLocalUsers()
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, data)
}

func DetailWindowsUser(c *gin.Context) {
	var req request.DetailWindowsUser
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}

	localUsers, err := wapi.ListLocalUsers()
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}

	var data *so.LocalUser
	for _, v := range localUsers {
		if v.Username == req.UserName {
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
