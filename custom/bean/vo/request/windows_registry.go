package request

type AddWindowsRegistry struct {
	FilePath string `json:"file_path" form:"file_path" binding:"required"`
}

type DeleteWindowsRegistry struct {
	FilePath string `json:"file_path" form:"file_path" binding:"required"`
}
