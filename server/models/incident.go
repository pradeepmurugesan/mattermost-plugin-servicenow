package models

// Incident model
type Incident struct {
	ShortDescription string `json:"short_description"`
	SysCreatedBy     string `json:"sys_created_by"`
	CreatedByID      string `json:"created_by"`
	Priority         string `json:"priority"`
	Impact           string `json:"impact"`
}
