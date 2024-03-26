package request

type AddWindowsUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	FullName string `json:"full_name" form:"full_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UpdateFullNameWindowsUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	FullName string `json:"full_name" form:"full_name" binding:"required"`
}

type UpdatePasswordWindowsUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type DeleteWindowsUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}

type DetailWindowsUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}
