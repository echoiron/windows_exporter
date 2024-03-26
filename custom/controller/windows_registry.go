package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	wapi "github.com/iamacarpet/go-win64api"
	"github.com/prometheus-community/windows_exporter/custom/bean/vo/request"
	"github.com/prometheus-community/windows_exporter/custom/resp"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
	"strings"
)

func AddWindowsRegistry(c *gin.Context) {
	var req request.AddWindowsRegistry
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}

	err := createRegistryEntry(req.FilePath)
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, nil)
}

func DeleteWindowsRegistry(c *gin.Context) {
	var req request.DeleteWindowsRegistry
	if err := c.ShouldBind(&req); err != nil {
		resp.SendJsonResponse(c, resp.InvalidParam.TranslateError(err), nil)
		return
	}
	err := deleteRegistryEntry(req.FilePath)
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, nil)
}

func ListWindowsRegistry(c *gin.Context) {
	data, err := wapi.ListLocalUsers()
	if err != nil {
		resp.SendJsonResponse(c, resp.FAIL.WithMsg(err.Error()), nil)
		return
	}
	resp.SendJsonResponse(c, resp.OK, data)
}

// 远程应用注册表项路径
var applicationsPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Terminal Server\TSAppAllowList\Applications`

func createRegistryEntry(filePath string) error {
	// 获取不带扩展名的文件名
	fileName := filepath.Base(filePath)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	keyPath := fmt.Sprintf("%s\\%s", applicationsPath, fileName)

	// 打开或创建上级目录
	if _, err := registry.OpenKey(registry.LOCAL_MACHINE, applicationsPath, registry.ALL_ACCESS); err != nil {
		parentKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, applicationsPath, registry.ALL_ACCESS)
		if err != nil {
			return fmt.Errorf("无法创建上级目录: %v", err)
		}
		defer parentKey.Close()
	}

	// 打开或创建注册表项
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("无法打开或创建注册表项: %v", err)
	}
	defer k.Close()

	// 设置键值对
	err = k.SetStringValue("Name", fileName)
	if err != nil {
		return fmt.Errorf("设置Name键值对时出错: %v", err)
	}

	err = k.SetStringValue("Path", filePath)
	if err != nil {
		return fmt.Errorf("设置Path键值对时出错: %v", err)
	}

	err = k.SetStringValue("VPath", filePath)
	if err != nil {
		return fmt.Errorf("设置VPath键值对时出错: %v", err)
	}

	err = k.SetStringValue("RequiredCommandLine", "")
	if err != nil {
		return fmt.Errorf("设置RequiredCommandLine键值对时出错: %v", err)
	}

	err = k.SetDWordValue("CommandLineSetting", 0x00000001)
	if err != nil {
		return fmt.Errorf("设置CommandLineSetting键值对时出错: %v", err)
	}

	err = k.SetStringValue("IconPath", filePath)
	if err != nil {
		return fmt.Errorf("设置IconPath键值对时出错: %v", err)
	}

	err = k.SetDWordValue("IconIndex", 0x00000000)
	if err != nil {
		return fmt.Errorf("设置IconIndex键值对时出错: %v", err)
	}

	err = k.SetDWordValue("ShowInTSWA", 0x00000000)
	if err != nil {
		return fmt.Errorf("设置ShowInTSWA键值对时出错: %v", err)
	}

	fmt.Println("注册表项已创建或已存在：" + keyPath)
	return nil
}

func deleteRegistryEntry(filePath string) error {
	// 获取不带扩展名的文件名
	fileName := filepath.Base(filePath)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	keyPath := fmt.Sprintf("%s\\%s", applicationsPath, fileName)

	// 删除注册表项
	err := registry.DeleteKey(registry.LOCAL_MACHINE, keyPath)
	if err != nil {
		return fmt.Errorf("无法删除注册表项: %v", err)
	}

	fmt.Println("注册表项已删除：" + keyPath)
	return nil
}
