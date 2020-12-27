package scrumProjectUseCase

import (
	"errors"
	"flowban/model"
	"flowban/module/scrumProject"
	"flowban/module/scrumProject/scrumProjectDTO"
)

type scrumProjectUseCase struct {
	scrumProjectRepo scrumProject.ScrumProjectRepository
}

//NewscrumProjectUseCase return implementation of ScrumProjectRepository Use Case
func NewScrumProjectRepositoryUseCase(scrumProjectRepo scrumProject.ScrumProjectRepository) scrumProject.ScrumProjectUseCase {
	return &scrumProjectUseCase{scrumProjectRepo: scrumProjectRepo}
}

func (u scrumProjectUseCase) GetProjectBoards(projectID int) (data *[]model.ScrumKanban, err error) {
	return u.scrumProjectRepo.GetBoardByProjectID(projectID)
}

func (u scrumProjectUseCase) MoveIssueToBacklog(issueID int) (*model.SprintIssue, error) {
	issue, err := u.scrumProjectRepo.GetIssueByID(issueID)
	if err != nil {
		return nil, err
	}
	issue.ScrumKanbanID = 0
	return u.scrumProjectRepo.UpdateIssue(*issue)
}

func (u scrumProjectUseCase) GetByUserID(userID int) (data *[]model.ScrumProject, err error) {
	return u.scrumProjectRepo.GetByUserID(userID)
}

func (u scrumProjectUseCase) MoveIssue(issueId int, boardID int) (*model.SprintIssue, error) {
	//get issue by id
	issueData, err := u.scrumProjectRepo.GetIssueByID(issueId)
	if err != nil {
		return nil, err
	}

	issueData.ScrumKanbanID = boardID

	return u.scrumProjectRepo.UpdateIssue(*issueData)
}

func (u scrumProjectUseCase) MoveIssueFromBacklog(issueID int) (*model.SprintIssue, error) {
	//get issue by id
	issueData, err := u.scrumProjectRepo.GetIssueByID(issueID)
	if err != nil {
		return nil, err
	}

	//get project by ID
	project, err := u.scrumProjectRepo.GetByID(issueData.ProjectID, false)
	if err != nil {
		return nil, err
	}

	//get board with name "To Do"
	board, err := u.scrumProjectRepo.GetBoardByProjectIDAndBoardName(project.ID, "To Do")
	if err != nil {
		return nil, err
	}

	//update issue
	issueData.ScrumKanbanID = board.ID
	return u.scrumProjectRepo.UpdateIssue(*issueData)
}

func (u scrumProjectUseCase) AddNewSprint(sprintData scrumProjectDTO.ReqCreateNewSprint) (*model.SprintSession, error) {
	//create new sprint and close current sprint if exist
	project, err := u.scrumProjectRepo.GetByID(sprintData.ScrumProjectID, true)
	if err != nil {
		return nil, err
	}

	sprintCount := 1
	sprintCount += project.CurrentSprint.SprintCount

	//build sprint
	sprintModel := model.SprintSession{
		Name:           sprintData.Name,
		Goal:           sprintData.Goal,
		StartDate:      sprintData.StartDate,
		EndDate:        sprintData.EndDate,
		ScrumProjectID: sprintData.ScrumProjectID,
		SprintCount:    sprintCount,
	}

	newSprint, err := u.scrumProjectRepo.InsertSprint(sprintModel)
	if err != nil {
		return nil, err
	}

	project.CurrentSprintID = newSprint.ID
	_, err = u.scrumProjectRepo.Update(*project)
	if err != nil {
		return nil, err
	}

	return newSprint, nil
}

func (u scrumProjectUseCase) AddNewIssue(issueData scrumProjectDTO.ReqCreateNewIssue) (*model.SprintIssue, error) {
	//get project data
	project, err := u.scrumProjectRepo.GetByID(issueData.ProjectID, true)
	if err != nil {
		return nil, err
	}

	if project.CurrentSprintID == 0 {
		return nil, errors.New("cannot add unstarted project, please create new sprint first")
	}

	newIssue := model.SprintIssue{
		Name:            issueData.Name,
		AssigneeID:      issueData.AssigneeID,
		ProjectID:       issueData.ProjectID,
		StoryPoint:      issueData.StoryPoint,
		SprintSessionID: project.CurrentSprintID,
	}

	savedIssue, err := u.scrumProjectRepo.InsertIssue(newIssue)
	return savedIssue, err
}

func (u scrumProjectUseCase) GetActiveIssueByProjectID(projectID int) (*[]model.SprintIssue, error) {
	//get project data
	project, err := u.scrumProjectRepo.GetByID(projectID, true)
	if err != nil {
		return nil, err
	}

	issues, err := u.scrumProjectRepo.GetByProjectIDAndSprintID(projectID, project.CurrentSprintID)
	return issues, err
}

func (u scrumProjectUseCase) GetAll() (allData *[]model.ScrumProject, err error) {
	return u.scrumProjectRepo.GetAll(true)
}

func (u scrumProjectUseCase) GetByID(dataId int) (data *model.ScrumProject, err error) {
	return u.scrumProjectRepo.GetByID(dataId, true)
}

func (u scrumProjectUseCase) AddNew(newData model.ScrumProject) (returnedNewData *model.ScrumProject, err error) {
	//Check Duplicate
	/*
		Block Code here
	*/

	newAddedData, err := u.scrumProjectRepo.Insert(newData)
	if err != nil {
		return nil, err
	}

	//add new default board
	projectBoards := []model.ScrumKanban{
		{
			ScrumProjectID: newAddedData.ID,
			BoardName:      "To Do",
		},
		{
			ScrumProjectID: newAddedData.ID,
			BoardName:      "In Progress",
		},
		{
			ScrumProjectID: newAddedData.ID,
			BoardName:      "Done",
		},
	}

	for _, kanban := range projectBoards {
		_, err = u.scrumProjectRepo.InsertNewBoard(kanban)
		if err != nil {
			return nil, err
		}
	}

	return newAddedData, err
}

func (u scrumProjectUseCase) Update(updatedData model.ScrumProject) (newUpdatedData *model.ScrumProject, err error) {
	//Check Duplicate
	/*
		Block Code here
	*/

	newUpdatedData, err = u.scrumProjectRepo.Update(updatedData)
	return newUpdatedData, err
}

func (u scrumProjectUseCase) RemoveByID(dataId int) error {
	errDelete := u.scrumProjectRepo.DeleteByID(dataId)
	return errDelete
}
