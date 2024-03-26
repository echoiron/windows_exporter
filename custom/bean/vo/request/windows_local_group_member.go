package request

type DetailWindowsLocalGroupMembers struct {
	GroupName string `json:"group_name" form:"group_name"`
	UserName  string `json:"user_name" form:"user_name" binding:"required"`
}

type ListWindowsLocalGroupMembers struct {
	GroupName string `json:"group_name" form:"group_name"`
}

type AddWindowsLocalGroupMembers struct {
	GroupName string   `json:"group_name" form:"group_name"`
	UserName  []string `json:"user_name" form:"user_name" binding:"required,min=1"`
}

type DeleteWindowsLocalGroupMembers struct {
	GroupName string   `json:"group_name" form:"group_name"`
	UserName  []string `json:"user_name" form:"user_name" binding:"required,min=1"`
}
