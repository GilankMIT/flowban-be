package model

import "flowban/helper/dbAudit"

type SprintIssue struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	AssigneeID      int         `json:"assignee_id"`
	Assignee        User        `json:"assignee"`
	ScrumKanbanID   int         `json:"scrum_kanban_id"`
	ScrumKanban     ScrumKanban `json:"scrum_kanban"`
	ProjectID       int         `json:"project_id"`
	SprintSessionID int         `json:"sprint_session_id"`
	StoryPoint      int         `json:"story_point"`
	dbAudit.DateAudit
	dbAudit.UserAudit
}
