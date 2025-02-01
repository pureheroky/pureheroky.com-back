package models

type UserRequest struct {
	UserName    string `json:"username"`
	UserEmail   string `json:"useremail"`
	UserMessage string `json:"usermessage"`
}
