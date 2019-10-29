package models

type Incident struct {
	ShortDescription string `json:"short_description"`
	SysCreatedBy     string `json:"sys_created_by"`
	CreatedById      string `json:"created_by"`
	Priority         string `json:"priority"`
	Impact           string `json:"impact"`
}
