package authUseCase

import (
	"errors"
	"flowban/helper/bcryptHelpers"
	"flowban/helper/jwthelper"
	"flowban/helper/sendgridMailHelper"
	"flowban/helper/stringBuilder"
	"flowban/model"
	"flowban/module/auth"
	"flowban/module/auth/authDTO"
	"flowban/module/user"
	"flowban/module/user/userUsecase"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

const (
	DefaultHashCost = 14
)

var (
	ErrInvalidCredential             = errors.New("authUseCase.authUseCase : invalid user credential")
	ErrUserNotFound                  = errors.New("authUseCase.authUseCase : user not found")
	ErrPasswordMismatch              = errors.New("authUseCase.authUseCase : password mismatch")
	ErrExpiredEmailVerificationToken = errors.New("authUseCase.authUseCase : token sent by email is expired")
	ErrInvalidEmailVerificationToken = errors.New("authUseCase.authUseCase : invalid email verification token")
	ErrUserIsInActive                = errors.New("authUsecase.authUseCase : user is inactive, please contact admin")
	ErrEmailAlreadyExist             = errors.New("email already exist")
)

type authUseCase struct {
	userRepo     user.UserRepository
	emailService sendgridMailHelper.SendGridMailService
}

func NewAuthUseCase(userRepo user.UserRepository, emailService sendgridMailHelper.SendGridMailService) auth.AuthUseCase {
	return &authUseCase{userRepo: userRepo,
		emailService: emailService}
}

func (a authUseCase) Register(reqRegister *authDTO.ReqRegister) (*authDTO.ResRegister, error) {
	//check if email already exist
	_, err := a.userRepo.GetByEmail(reqRegister.Email, false)
	if err == nil {
		return nil, ErrEmailAlreadyExist
	}

	//generate model
	userModel := model.User{
		Email:     reqRegister.Email,
		FirstName: reqRegister.FirstName,
		LastName:  reqRegister.LastName,
		UserRole: []model.Role{
			{Name: "user"},
		},
	}

	//generate bcrypt for password
	passwordHashed, _ := bcryptHelpers.GenerateBcryptCustomCost(reqRegister.Password, 14)
	userModel.Password = passwordHashed

	newUser, err := a.userRepo.Insert(userModel)
	if err != nil {
		return nil, err
	}

	//generate token from user
	jwtExpirationDurationDayString := os.Getenv("jwt.expirationDurationDay")
	var jwtExpirationDurationDay int
	jwtExpirationDurationDay, err = strconv.Atoi(jwtExpirationDurationDayString)
	if err != nil {
		return nil, err
	}

	// Conversion to seconds
	jwtExpiredAt := time.Now().Unix() + int64(jwtExpirationDurationDay*3600*24)

	//Combine Role into String
	//Separated with "|"
	userRoles := newUser.UserRole
	var userRolesString string
	for i, role := range userRoles {
		if i < len(userRoles)-1 {
			userRolesString = userRolesString + role.Name + "|"
		} else {
			userRolesString = userRolesString + role.Name
		}
	}

	userClaims := jwthelper.CustomClaims{Role: userRolesString, Id: newUser.ID, ExpiresAt: jwtExpiredAt}
	jwtToken, err := jwthelper.NewWithClaims(userClaims)
	if err != nil {
		return nil, err
	}

	res := authDTO.ResRegister{Token: jwtToken}
	return &res, nil
}

func (a authUseCase) Login(reqLogin *authDTO.ReqLoginEmailPass) (*authDTO.ResLogin, error) {
	userData, err := a.userRepo.GetByEmail(reqLogin.Email, true)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}

	//check user status, if inactive, return 401
	if userData.Status == userUsecase.UserStatusInactive {
		return nil, ErrUserIsInActive
	}

	err = bcryptHelpers.CompareBcrypt(userData.Password, reqLogin.Password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrInvalidCredential
		} else {
			return nil, err
		}
	}

	jwtExpirationDurationDayString := os.Getenv("jwt.expirationDurationDay")
	var jwtExpirationDurationDay int
	jwtExpirationDurationDay, err = strconv.Atoi(jwtExpirationDurationDayString)
	if err != nil {
		return nil, err
	}

	// Conversion to seconds
	jwtExpiredAt := time.Now().Unix() + int64(jwtExpirationDurationDay*3600*24)

	//Combine Role into String
	//Separated with "|"
	userRoles := userData.UserRole
	var userRolesString string
	for i, role := range userRoles {
		if i < len(userRoles)-1 {
			userRolesString = userRolesString + role.Name + "|"
		} else {
			userRolesString = userRolesString + role.Name
		}
	}

	userClaims := jwthelper.CustomClaims{Role: userRolesString, Id: userData.ID, ExpiresAt: jwtExpiredAt}
	jwtToken, err := jwthelper.NewWithClaims(userClaims)
	if err != nil {
		return nil, err
	}

	loginRes := authDTO.ResLogin{User: userData, Token: jwtToken}
	return &loginRes, nil
}

func (a authUseCase) ChangePassword(userId int, prevPassword, newPassword string) error {
	//find user by user id
	userData, err := a.userRepo.GetByID(userId, true)
	if err != nil {
		return ErrUserNotFound
	}

	//validate password first
	err = bcryptHelpers.CompareBcrypt(userData.Password, prevPassword)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrPasswordMismatch
		}
		return err
	}

	//create new password to user
	hashedPassword, err := bcryptHelpers.GenerateBcryptCustomCost(newPassword, DefaultHashCost)
	if err != nil {
		return err
	}

	userData.Password = hashedPassword

	//save new password to db
	_, err = a.userRepo.Update(*userData)
	if err != nil {
		return err
	}

	return nil
}

func (a authUseCase) SendPasswordResetByEmail(email string) (message string, err error) {
	//retrieve user data by id
	userData, err := a.userRepo.GetByEmail(email, true)
	if err != nil {
		return "", ErrUserNotFound
	}

	//create random token
	//with structure : UserID_RandomStringAlphaNum
	const (
		emailVerificationTokenLength      = 72
		emailVerificationTokenExpDuration = time.Hour * 24 * 5 //5 days
	)

	emailVerificationToken := stringBuilder.GenerateRandom(emailVerificationTokenLength)
	userData.EmailVerificationToken = strconv.Itoa(userData.ID) + "_" + emailVerificationToken
	userData.EmailVerificationTokenExp = time.Now().Add(emailVerificationTokenExpDuration).Unix()

	//save token to repo
	_, err = a.userRepo.Update(*userData)

	//send token to email
	baseUrl := os.Getenv("base_url")
	emailPayload := sendgridMailHelper.Email{
		Receivers:   []string{userData.Email},
		Subject:     "Email reset verification",
		ContentType: "text/html",
		Content: fmt.Sprintf("This is the link to reset your password <br> "+
			"<a href='%s/api/v1/auth/view/reset-password?token=%s'>Click here to reset</a>",
			baseUrl, userData.EmailVerificationToken),
	}
	noreplySender := os.Getenv("sendgrid.sender_mail")
	err = a.emailService.SendMail(emailPayload, noreplySender)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (a authUseCase) ResetPasswordByToken(userId int, emailVerificationToken, newPassword string) error {
	//retrieve user by ID
	//retrieve user data by id
	userData, err := a.userRepo.GetByID(userId, true)
	if err != nil {
		return ErrUserNotFound
	}

	//compare token against the saved one in DB
	if !(userData.EmailVerificationTokenExp > time.Now().Unix()) {
		return ErrExpiredEmailVerificationToken
	}

	if userData.EmailVerificationToken != emailVerificationToken {
		return ErrInvalidEmailVerificationToken
	}

	//token is valid, reset user password and make token expired
	hashedPassword, err := bcryptHelpers.GenerateBcryptCustomCost(newPassword, DefaultHashCost)
	if err != nil {
		return err
	}
	userData.Password = hashedPassword
	userData.EmailVerificationTokenExp = time.Now().Add(-time.Second * 1).Unix()
	//save new password to db
	_, err = a.userRepo.Update(*userData)
	if err != nil {
		return err
	}

	return nil
}
