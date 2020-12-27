package user

import "flowban/model"

type UserRepository interface {
	GetterRepository
	InserterRepository
	UpdaterRepository
	DeleterRepository
}

//GetterRepository define the repository for data retrieval scenario
type GetterRepository interface {
	GetAll(autoPreload bool) (*[]model.User, error)
	GetAllRole() (*[]model.Role, error)
	GetAllWhereActive(status int, autoPreload bool) (*[]model.User, error)
	GetByID(dataId int, autoPreload bool) (*model.User, error)
	GetByEmail(email string, autoPreload bool) (*model.User, error)
}

//InserterRepository define the repository for data insertion scenario
type InserterRepository interface {
	Insert(data model.User) (*model.User, error)
}

//UpdaterRepository define the repository for data modification scenario
type UpdaterRepository interface {
	SetStatusForUserID(userId int, status int) (*model.User, error)
	Update(data model.User) (*model.User, error)
}

//DeleterRepository define the repository for data removal scenario
type DeleterRepository interface {
	DeleteByID(dataID int) error
}
