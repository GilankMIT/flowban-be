package scrumProject

import "flowban/model"

type ScrumProjectRepository interface {
	GetterRepository
	InserterRepository
	UpdaterRepository
	DeleterRepository
}

//GetterRepository define the repository for data retrieval scenario
type GetterRepository interface {
	GetAll(autoPreload bool) (*[]model.ScrumProject, error)
	GetByID(dataId int, autoPreload bool) (*model.ScrumProject, error)
	GetByUserID(userId int) (*[]model.ScrumProject, error)
	GetByProjectIDAndSprintID(projectId, sprintId int) (*[]model.SprintIssue, error)
	GetIssueByID(id int) (*model.SprintIssue, error)
	GetBoardByProjectIDAndBoardName(id int, name string) (*model.ScrumKanban, error)
	GetBoardByProjectID(projectID int) (*[]model.ScrumKanban, error)
}

//InserterRepository define the repository for data insertion scenario
type InserterRepository interface {
	Insert(data model.ScrumProject) (*model.ScrumProject, error)
	InsertSprint(data model.SprintSession) (*model.SprintSession, error)
	InsertIssue(data model.SprintIssue) (*model.SprintIssue, error)
	InsertNewBoard(data model.ScrumKanban) (*model.ScrumKanban, error)
}

//UpdaterRepository define the repository for data modification scenario
type UpdaterRepository interface {
	Update(data model.ScrumProject) (*model.ScrumProject, error)
	UpdateIssue(data model.SprintIssue) (*model.SprintIssue, error)
}

//DeleterRepository define the repository for data removal scenario
type DeleterRepository interface {
	DeleteByID(dataID int) error
}
