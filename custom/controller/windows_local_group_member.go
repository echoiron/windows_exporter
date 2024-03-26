package controller

import (
	"github.com/gin-gonic/gin"
	wapi "github.com/iamacarpet/go-win64api"
	so "github.com/iamacarpet/go-win64api/shared"
	"github.com/prometheus-community/windows_exporter/custom/bean/vo/request"
	"github.com/prometheus-community/windows_exporter/custom/resp"
	"strings"
)

var defaultGroupName = "Remote Desktop Users"

func AddWindowsLocalGroupMembers(c *gin.Context) {
	var req request.AddWindowsLocalGroupMembers
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}

	if req.GroupName == "" {
		req.GroupName = defaultGroupName
	}
	ok, err := wapi.LocalGroupAddMembers(req.GroupName, req.UserName)
	commonReturn(c, err, ok)
}

func DeleteWindowsLocalGroupMembers(c *gin.Context) {
	var req request.DeleteWindowsLocalGroupMembers
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	if req.GroupName == "" {
		req.GroupName = defaultGroupName
	}
	ok, err := wapi.LocalGroupDelMembers(req.GroupName, req.UserName)
	commonReturn(c, err, ok)
}

func ListWindowsLocalGroupMembers(c *gin.Context) {
	var req request.ListWindowsLocalGroupMembers
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}

	if req.GroupName == "" {
		req.GroupName = defaultGroupName
	}
	data, err := wapi.LocalGroupGetMembers(req.GroupName)
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, data)
}

func DetailWindowsLocalGroupMembers(c *gin.Context) {
	var req request.DetailWindowsLocalGroupMembers
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}

	if req.GroupName == "" {
		req.GroupName = defaultGroupName
	}
	members, err := wapi.LocalGroupGetMembers(req.GroupName)
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}

	var data *so.LocalGroupMember
	//  "name": "DESKTOP-T0FNQON\xiaowang",
	for _, v := range members {
		splits := strings.Split(v.Name, `\`)
		if len(splits) >= 2 {
			if splits[1] == req.UserName {
				data = &v
				break
			}
		}
	}

	if data == nil {
		resp.SendJsonResponse(c, resp.NotFound.WithMsg("数据不存在"), nil)
		return
	}

	resp.SendJsonResponse(c, resp.OK, data)
}
