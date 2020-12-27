package scrumProjectDTO

type ReqAddNewProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Acronym     string `json:"acronym"`
}

type ReqCreateNewSprint struct {
	Name           string `json:"name"`
	Goal           string `json:"goal"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	ScrumProjectID int    `json:"scrum_project_id"`
}

type ReqCreateNewIssue struct {
	Name       string `json:"name"`
	AssigneeID int    `json:"assignee_id"`
	ProjectID  int    `json:"project_id"`
	StoryPoint int    `json:"story_point"`
}
