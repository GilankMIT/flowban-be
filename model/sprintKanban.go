package model

import "flowban/helper/dbAudit"

type ScrumKanban struct {
	ID             int    `json:"id"`
	ScrumProjectID int    `json:"scrum_project_id"`
	BoardName      string `json:"board_name"`
	dbAudit.UserAudit
	dbAudit.DateAudit
}
