package response

import "github.com/prometheus-community/windows_exporter/custom/bean/entity"

type DataSourceItem struct {
	entity.WindowsUser
}

type DataSourceList struct {
	TotalCount int64            `json:"total_count"`
	List       []DataSourceItem `json:"list"`
}
