package models

type UserData struct {
	Username string `json:"username" bson:"username"`
	Age      int    `json:"age"      bson:"age"`
	Status   string `json:"status"   bson:"status"`
	Avatar   string `json:"avatar"   bson:"avatar"`
}

type UserResponse struct {
	Data   UserData `json:"data"`
	Status int      `json:"status"`
}
