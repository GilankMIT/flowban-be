package userRepository

import (
	"flowban/model"
	"flowban/module/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepo struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) user.UserRepository {
	return &userRepo{conn: conn}
}

func (u userRepo) GetAllWhereActive(status int, autoPreload bool) (*[]model.User, error) {
	var dataList []model.User

	db := u.conn
	//preload check
	if autoPreload {
		db = db.Preload(clause.Associations)
	}

	db = db.Find(&dataList, "status=?", status)
	return &dataList, db.Error
}

//GetAll retrieve all data from DB
func (u userRepo) GetAll(autoPreload bool) (*[]model.User, error) {
	var dataList []model.User

	db := u.conn
	//preload check
	if autoPreload {
		db = db.Preload(clause.Associations)
	}

	db = db.Find(&dataList)
	return &dataList, db.Error
}

//GetByID retrieve data by ID from DB
func (u userRepo) GetByID(dataId int, autoPreload bool) (*model.User, error) {
	var dataList model.User

	db := u.conn
	//preload check
	if autoPreload {
		db = db.Preload(clause.Associations)
	}

	db = db.First(&dataList, "id=?", dataId)
	return &dataList, db.Error
}

//GetByID retrieve data by Email from DB
func (u userRepo) GetByEmail(email string, autoPreload bool) (*model.User, error) {
	var dataList model.User

	db := u.conn
	//preload check
	if autoPreload {
		db = db.Preload(clause.Associations)
	}

	db = db.First(&dataList, "email=?", email)
	return &dataList, db.Error
}

//Insert add new data to DB
func (u userRepo) Insert(data model.User) (*model.User, error) {
	db := u.conn.Create(&data)
	return &data, db.Error
}

//Update modify existing data from DB
func (u userRepo) Update(data model.User) (*model.User, error) {
	db := u.conn.Save(&data)
	return &data, db.Error
}

func (u userRepo) SetStatusForUserID(userId int, status int) (*model.User, error) {
	db := u.conn.Model(&model.User{}).
		Where("id=?", userId).
		Update("status", status)
	if db.Error != nil {
		return nil, db.Error
	}

	var userData model.User
	db = u.conn.First(&userData, "id=?", userId)

	return &userData, db.Error
}

//DeleteByID remove data from DB
func (u userRepo) DeleteByID(dataID int) error {
	db := u.conn.Delete(&model.User{}, "id=?", dataID)
	return db.Error
}

func (u userRepo) GetAllRole() (*[]model.Role, error) {
	var roleList []model.Role
	db := u.conn.Find(&roleList)
	return &roleList, db.Error
}
