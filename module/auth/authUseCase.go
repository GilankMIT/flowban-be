package auth

import "flowban/module/auth/authDTO"

type AuthUseCase interface {
	Login(reqLogin *authDTO.ReqLoginEmailPass) (*authDTO.ResLogin, error)
	Register(reqRegister *authDTO.ReqRegister) (*authDTO.ResRegister, error)
	ChangePassword(userId int, prevPassword, newPassword string) error
	SendPasswordResetByEmail(email string) (message string, err error)
	ResetPasswordByToken(userId int, emailVerificationToken, newPassword string) error
}
