package entity

type WindowsUser struct {
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}
