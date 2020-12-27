package model

import "flowban/helper/dbAudit"

type ScrumKanban struct {
	ID             int          `json:"id"`
	ScrumProjectID int          `json:"scrum_project_id"`
	ScrumProject   ScrumProject `json:"scrum_project"`
	BoardName      string
	dbAudit.UserAudit
	dbAudit.DateAudit
}
