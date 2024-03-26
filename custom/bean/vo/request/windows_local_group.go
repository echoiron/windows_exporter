package request

type DetailWindowsLocalGroup struct {
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
}

type ListWindowsLocalGroup struct {
	GroupName string `json:"group_name" form:"group_name"`
}

type AddWindowsLocalGroup struct {
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
	Comment   string `json:"comment" form:"comment" binding:"required"`
}

type DeleteWindowsLocalGroup struct {
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
}
