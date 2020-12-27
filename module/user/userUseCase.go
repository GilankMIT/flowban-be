package user

import (
	"flowban/model"
	"flowban/module/user/userDTO"
)

type UserUseCase interface {
	//GetAll return all data of User from repo
	GetAll() (allData *[]model.User, err error)

	//GetByID return of User by ID from repo
	GetByID(dataId int) (data *model.User, err error)

	//AddNew add new data of User to repo
	AddNew(userData userDTO.ReqAddNewUser) (returnedNewData *model.User, err error)

	//Update modify existing User from repo
	Update(updatedDate model.User) (updatedData *model.User, err error)

	UpdateCredentialData(userId int, newPassword, newFirstName, newLastName string) (*model.User, error)

	//RemoveByID remove data User from repo
	RemoveByID(dataId int) error

	GetAllRole() (*[]model.Role, error)
}
