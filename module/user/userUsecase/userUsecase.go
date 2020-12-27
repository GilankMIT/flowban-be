package userUsecase

import (
	"errors"
	"flowban/helper/bcryptHelpers"
	"flowban/helper/sendgridMailHelper"
	"flowban/helper/stringBuilder"
	"flowban/model"
	"flowban/module/user"
	"flowban/module/user/userDTO"
	"fmt"
	"gorm.io/gorm"
	"os"
)

var (
	ErrEmailAlreadyExist = errors.New("userUC.userUsecase : email already exist")
	ErrUserNotFound      = errors.New("userUC.userUsecase : user not found")
)

const (
	UserStatusActive = iota + 1
	UserStatusInactive
)

type userUseCase struct {
	userRepo    user.UserRepository
	mailService sendgridMailHelper.SendGridMailService
}

//NewUserUseCase return implementation of User Use Case
func NewUserRepositoryUseCase(userRepo user.UserRepository, mailService sendgridMailHelper.SendGridMailService) user.UserUseCase {
	return &userUseCase{
		userRepo:    userRepo,
		mailService: mailService,
	}
}

func (u userUseCase) AddNew(userData userDTO.ReqAddNewUser) (returnedNewData *model.User, err error) {
	//duplication check
	_, err = u.userRepo.GetByEmail(userData.Email, false)
	if err == nil {
		return nil, ErrEmailAlreadyExist
	}

	//generate new password
	newPassword := stringBuilder.GenerateRandomCustom(8)
	userData.Password = newPassword

	emailPayload := sendgridMailHelper.Email{
		Receivers:   []string{userData.Email},
		Subject:     "Akun UPK Report",
		ContentType: "text/html",
		Content: fmt.Sprintf("Akun anda telah dibuat dengan email <b>" + userData.Email +
			"</b>. <br> Silakan login dengan password : <b>" + newPassword + "</b>"),
	}
	noreplySender := os.Getenv("sendgrid.sender_mail")
	err = u.mailService.SendMail(emailPayload, noreplySender)
	if err != nil {
		return nil, err
	}

	//hash bcrypt
	hashedPassword, _ := bcryptHelpers.GenerateBcryptCustomCost(userData.Password, 14)

	//create role
	roles := make([]model.Role, 0)
	for _, roleNode := range userData.Role {
		roles = append(roles, model.Role{
			ID: roleNode,
		})
	}

	//build model
	userDataModel := model.User{
		Password:  hashedPassword,
		Email:     userData.Email,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		UserRole:  roles,
		Status:    UserStatusActive,
	}

	return u.userRepo.Insert(userDataModel)
}

func (u userUseCase) UpdateCredentialData(userId int, newPassword, newFirstName, newLastName string) (*model.User, error) {
	//get user
	existingUser, err := u.userRepo.GetByID(userId, true)
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	//check if user change password
	if newPassword != "" {
		hashedNewPassword, _ := bcryptHelpers.GenerateBcryptCustomCost(newPassword, 14)
		existingUser.Password = hashedNewPassword
	}

	existingUser.FirstName = newFirstName
	existingUser.LastName = newLastName

	return u.userRepo.Update(*existingUser)
}

func (u userUseCase) GetAll() (allData *[]model.User, err error) {
	return u.userRepo.GetAllWhereActive(UserStatusActive, true)
}

func (u userUseCase) GetByID(dataId int) (data *model.User, err error) {
	return u.userRepo.GetByID(dataId, true)
}

func (u userUseCase) Update(updatedData model.User) (newUpdatedData *model.User, err error) {
	//Check Duplicate
	/*
		Block Code here
	*/

	newUpdatedData, err = u.userRepo.Update(updatedData)
	return newUpdatedData, err
}

func (u userUseCase) RemoveByID(dataId int) error {
	//check if user exit
	_, err := u.userRepo.GetByID(dataId, false)
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return ErrUserNotFound
		}
		return err
	}

	_, err = u.userRepo.SetStatusForUserID(dataId, UserStatusInactive)
	return err
}

func (u userUseCase) GetAllRole() (*[]model.Role, error) {
	return u.userRepo.GetAllRole()
}
