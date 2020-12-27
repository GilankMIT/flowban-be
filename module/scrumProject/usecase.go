package scrumProject

import (
	"flowban/model"
	"flowban/module/scrumProject/scrumProjectDTO"
)

type ScrumProjectUseCase interface {
	//GetAll return all data of ScrumProject from repo
	GetAll() (allData *[]model.ScrumProject, err error)

	//GetByID return of ScrumProject by ID from repo
	GetByID(dataId int) (data *model.ScrumProject, err error)

	GetByUserID(userID int) (data *[]model.ScrumProject, err error)

	GetProjectBoards(projectID int) (data *[]model.ScrumKanban, err error)

	//AddNew add new data of ScrumProject to repo
	AddNew(newData model.ScrumProject) (returnedNewData *model.ScrumProject, err error)

	AddNewSprint(sprintData scrumProjectDTO.ReqCreateNewSprint) (*model.SprintSession, error)

	AddNewIssue(issueData scrumProjectDTO.ReqCreateNewIssue) (*model.SprintIssue, error)

	GetActiveIssueByProjectID(projectID int) (*[]model.SprintIssue, error)

	MoveIssue(issueId int, boardID int) (*model.SprintIssue, error)

	MoveIssueFromBacklog(issueID int) (*model.SprintIssue, error)

	MoveIssueToBacklog(issueID int) (*model.SprintIssue, error)

	//Update modify existing ScrumProject from repo
	Update(updatedDate model.ScrumProject) (updatedData *model.ScrumProject, err error)

	//RemoveByID remove data ScrumProject from repo
	RemoveByID(dataId int) error
}
