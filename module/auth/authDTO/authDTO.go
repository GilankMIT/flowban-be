package authDTO

import "flowban/model"

//ReqLoginEmailPass - Object for mapping login request
type ReqLoginEmailPass struct {
	Email    string `json:"email" binding:"required" example:"kim@kecjambe.go.id"`
	Password string `json:"password" binding:"required" exampel:"kimwilde12345"`
}

//ResLogin - Object for mapping login response
type ResLogin struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

type ReqUpdatePassword struct {
	PreviousPassword string `json:"previous_password" binding:"required" example:"oldpass0987TiarGanteng"`
	NewPassword      string `json:"new_password" binding:"required" example:"newpass12345TiarMending"`
}

type ResUpdatePassword struct {
	Mesage string `json:"mesage" example:"password successfully changed"`
}

type ReqResetPasswordByEmail struct {
	Email string `json:"email" binding:"required" example:"maria@loophole.com"`
}

type ResResetPasswordByEmail struct {
	Message string `json:"message" example:"email reset request successfully sent to [EMAIL]"`
}

//ReqGetUserByJWT - Object for mapping Get User by JWT request
type ReqGetUserByJWT struct {
	Token string `json:"token"`
	Roles string `json:"roles"`
}

//ResGetUserByJWT - Object for mapping Get User by JWT response
type ResGetUserByJWT struct {
	User *model.User `json:"user"`
}

type ReqRegister struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ResRegister struct {
	Token string `json:"token"`
}
