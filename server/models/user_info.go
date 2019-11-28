package models

type result struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

// UserInfo user details from service now
type UserInfo struct {
	Result result `json:"result"`
}
