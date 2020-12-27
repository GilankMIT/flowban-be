package model

import (
	"flowban/helper/dbAudit"
)

type SprintSession struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Goal           string `json:"goal"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	ScrumProjectID int    `json:"scrum_project_id"`
	SprintCount    int    `json:"sprint_count"`
	dbAudit.UserAudit
	dbAudit.DateAudit
}
