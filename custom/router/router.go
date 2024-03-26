package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus-community/windows_exporter/custom/controller"
)

func RegisterRoutes(e *gin.Engine) {
	group := e.Group("windows")
	{
		user := group.Group("user")
		{
			user.GET("", controller.DetailWindowsUser)
			user.GET("list", controller.ListWindowsUser)
			user.POST("", controller.AddWindowsUser)
			user.PUT("full_name", controller.UpdateFullNameWindowsUser)
			user.PUT("password", controller.UpdatePasswordWindowsUser)
			user.DELETE("", controller.DeleteWindowsUser)
		}

		localGroup := group.Group("local_group")
		{
			localGroup.GET("", controller.DetailWindowsLocalGroup)
			localGroup.GET("list", controller.ListWindowsLocalGroup)
			localGroup.POST("", controller.AddWindowsLocalGroup)
			localGroup.DELETE("", controller.DeleteWindowsLocalGroup)

			localGroup.GET("members", controller.DetailWindowsLocalGroupMembers)
			localGroup.GET("members/list", controller.ListWindowsLocalGroupMembers)
			localGroup.POST("members", controller.AddWindowsLocalGroupMembers)
			localGroup.DELETE("members", controller.DeleteWindowsLocalGroupMembers)
		}

		adminGroup := group.Group("admin")
		{
			adminGroup.POST("set_admin", controller.SetAdmin)
			adminGroup.POST("revoke_admin", controller.RevokeAdmin)
		}

		registry := group.Group("registry")
		{
			registry.POST("", controller.AddWindowsRegistry)
			registry.DELETE("", controller.DeleteWindowsRegistry)
		}

	}
}
