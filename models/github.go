package models

import "time"

type GithubRepo struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

type GithubCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

type CommitInfo struct {
	ProjectName string    `json:"project_name"`
	Branch      string    `json:"branch"`
	Date        time.Time `json:"date"`
	Message     string    `json:"message"`
}
