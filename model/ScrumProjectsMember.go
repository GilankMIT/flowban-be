package model

type ScrumProjectMember struct {
	ID             int  `json:"id"`
	ScrumProjectID int  `json:"scrum_project_id"`
	UserID         int  `json:"user_id"`
	User           User `json:"user"`
}
