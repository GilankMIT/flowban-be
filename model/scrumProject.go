package model

import "flowban/helper/dbAudit"

type ScrumProject struct {
	ID              int                  `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	UserID          int                  `json:"scrum_master_id"`
	User            User                 `json:"scrum_master"`
	Members         []ScrumProjectMember `json:"members"`
	CurrentSprintID int                  `json:"current_sprint_id"`
	CurrentSprint   SprintSession        `json:"current_sprint"`
	ImageURL        string               `json:"image_url"`
	Acronym         string               `json:"acronym"`
	dbAudit.DateAudit
	dbAudit.UserAudit
}
