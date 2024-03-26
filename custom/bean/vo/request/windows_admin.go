package request

type SetAdmin struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}

type RevokeAdmin struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}
