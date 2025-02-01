package models

type Project struct {
	Name              string `json:"name"               bson:"name"`
	Description       string `json:"description"        bson:"description"`
	Date              string `json:"date"               bson:"date"`
	Image             string `json:"image"              bson:"image"`
	GitLink           string `json:"git_link"           bson:"git_link"`
	WebLink           string `json:"web_link"           bson:"web_link"`
	DevelopmentStatus string `json:"development_status" bson:"development_status"`
}

type ProjectResponse struct {
	Projects []Project `json:"projects"`
	Status   int       `json:"status"`
}
