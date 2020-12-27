package scrumProjectRepository

import (
	"flowban/model"
	"flowban/module/scrumProject"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type scrumProjectRepo struct {
	conn *gorm.DB
}

func NewScrumProjectRepository(conn *gorm.DB) scrumProject.ScrumProjectRepository {
	return &scrumProjectRepo{conn: conn}
}

func (u scrumProjectRepo) UpdateIssue(data model.SprintIssue) (*model.SprintIssue, error) {
	db := u.conn.Save(&data)
	return &data, db.Error
}

func (u scrumProjectRepo) GetIssueByID(id int) (*model.SprintIssue, error) {
	var issue model.SprintIssue
	db := u.conn.First(&issue, "id=?", id)
	return &issue, db.Error
}

func (u scrumProjectRepo) GetBoardByProjectIDAndBoardName(id int, name string) (*model.ScrumKanban, error) {
	var board model.ScrumKanban
	db := u.conn.First(&board, "scrum_project_id=? and board_name=?", id, name)
	return &board, db.Error
}

func (u scrumProjectRepo) GetBoardByProjectID(projectID int) (*[]model.ScrumKanban, error) {
	var board []model.ScrumKanban
	db := u.conn.Find(&board, "scrum_project_id=?", projectID)
	return &board, db.Error
}

func (u scrumProjectRepo) InsertNewBoard(data model.ScrumKanban) (*model.ScrumKanban, error) {
	db := u.conn.Create(&data)
	return &data, db.Error
}

func (u scrumProjectRepo) InsertSprint(data model.SprintSession) (*model.SprintSession, error) {
	db := u.conn.Create(&data)
	return &data, db.Error
}

func (u scrumProjectRepo) InsertIssue(data model.SprintIssue) (*model.SprintIssue, error) {
	db := u.conn.Create(&data)
	return &data, db.Error
}

func (u scrumProjectRepo) GetByUserID(userId int) (*[]model.ScrumProject, error) {
	var scrumProjects []model.ScrumProject
	db := u.conn.Preload(clause.Associations).Find(&scrumProjects, "user_id=?", userId)
	return &scrumProjects, db.Error
}

func (u scrumProjectRepo) GetByProjectIDAndSprintID(projectId, sprintId int) (*[]model.SprintIssue, error) {
	var sprintIssues []model.SprintIssue
	db := u.conn.Preload(clause.Associations).
		Find(&sprintIssues, "project_id = ? AND sprint_session_id = ?", projectId, sprintId)
	return &sprintIssues, db.Error
}

//GetAll retrieve all data from DB
func (u scrumProjectRepo) GetAll(autoPreload bool) (*[]model.ScrumProject, error) {
	var dataList []model.ScrumProject

	db := u.conn
	//preload check
	if autoPreload {
		db = db.Preload(clause.Associations)
	}

	db = db.Find(&dataList)
	return &dataList, db.Error
}

//GetByID retrieve data by ID from DB
func (u scrumProjectRepo) GetByID(dataId int, autoPreload bool) (*model.ScrumProject, error) {
	var dataList model.ScrumProject

	db := u.conn
	//preload check
	if autoPreload {
		db = db.Preload(clause.Associations)
	}

	db = db.First(&dataList, "id=?", dataId)
	return &dataList, db.Error
}

//Insert add new data to DB
func (u scrumProjectRepo) Insert(data model.ScrumProject) (*model.ScrumProject, error) {
	db := u.conn.Create(&data)
	return &data, db.Error
}

//Update modify existing data from DB
func (u scrumProjectRepo) Update(data model.ScrumProject) (*model.ScrumProject, error) {
	db := u.conn.Save(&data)
	return &data, db.Error
}

//DeleteByID remove data from DB
func (u scrumProjectRepo) DeleteByID(dataID int) error {
	db := u.conn.Delete(&model.ScrumProject{}, "id=?", dataID)
	return db.Error
}
