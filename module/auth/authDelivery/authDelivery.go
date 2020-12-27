package authDelivery

import (
	"flowban/helper/jsonHttpResponse"
	"flowban/helper/requestvalidationerror"
	"flowban/helper/routeHelpers"
	"flowban/middleware/authmiddleware"
	"flowban/module/auth"
	"flowban/module/auth/authDTO"
	"flowban/module/auth/authUseCase"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type authDelivery struct {
	authUC auth.AuthUseCase
}

func NewAuthHTTPDelivery(r *gin.Engine, authUC auth.AuthUseCase, authMiddleware authmiddleware.AuthMiddleware) {
	handler := authDelivery{authUC: authUC}

	publicAccess := r.Group("/api/v1/auth")
	{
		publicAccess.POST("login", handler.login)
		publicAccess.POST("register", handler.register)
		publicAccess.POST("reset-password/verify-reset", handler.verifyReset)
		publicAccess.POST("reset-password/send-mail", handler.resetPasswordByEmailSend)
	}

	authorizedUserAccess := r.Group("/api/v1/auth", authMiddleware.AuthorizeJWTWithUserContext("admin"))
	{
		authorizedUserAccess.POST("change-password", handler.changePassword)
	}

	//r.LoadHTMLGlob("views/**/*")
	//r.GET("/api/v1/auth/view/reset-password", func(c *gin.Context) {
	//	token := c.Query("token")
	//	baseUrl := os.Getenv("base_url")
	//	appName := os.Getenv("application.name")
	//	c.HTML(200, "auth/reset-password.html", gin.H{
	//		"baseUrl": baseUrl,
	//		"token":   token,
	//		"appName": appName,
	//	})
	//})
}

// Authentication - Login godoc
// @Summary Login - Get User data and token by email & password
// @Description Return JSON Web Token to authorize routes by using email and password
// @ID authentication-login
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param applicationData body authDTO.ReqLoginEmailPass true "Login data"
// @Success 200 {object} jsonHttpResponse.SuccessResponse{data=authDTO.ResLogin}
// @Failure 400 {object} jsonHttpResponse.FailedBadRequestResponse
// @Failure 401 {object} jsonHttpResponse.FailedUnauthorizedResponse
// @Failure 500 {object} jsonHttpResponse.FailedInternalServerErrorResponse
// @Router /api/v1/auth/login [post]
func (handler *authDelivery) login(c *gin.Context) {
	var request authDTO.ReqLoginEmailPass
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			jsonHttpResponse.NewFailedBadRequestResponse(c, validations)
			return
		}

		jsonHttpResponse.NewFailedBadRequestResponse(c, errBind.Error())
		return
	}

	res, err := handler.authUC.Login(&request)
	if err != nil {
		if err == authUseCase.ErrPasswordMismatch ||
			err == authUseCase.ErrInvalidCredential ||
			err == authUseCase.ErrUserIsInActive {
			jsonHttpResponse.NewFailedUnauthorizedResponse(c, err.Error())
			return
		}

		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, res)
	return
}

func (handler *authDelivery) register(c *gin.Context) {
	var request authDTO.ReqRegister
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			jsonHttpResponse.NewFailedBadRequestResponse(c, validations)
			return
		}

		jsonHttpResponse.NewFailedBadRequestResponse(c, errBind.Error())
		return
	}

	res, err := handler.authUC.Register(&request)
	if err != nil {
		if err == authUseCase.ErrPasswordMismatch ||
			err == authUseCase.ErrInvalidCredential ||
			err == authUseCase.ErrUserIsInActive {
			jsonHttpResponse.NewFailedUnauthorizedResponse(c, err.Error())
			return
		}

		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, res)
	return
}

// Authentication - Change password godoc
// @Summary Login - Change user password
// @Description Change user password. Validated by previous password
// @ID authentication-change-password
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param applicationData body authDTO.ReqUpdatePassword true "Change password data"
// @Success 200 {object} jsonHttpResponse.SuccessResponse{data=authDTO.ResUpdatePassword}
// @Failure 400 {object} jsonHttpResponse.FailedBadRequestResponse
// @Failure 401 {object} jsonHttpResponse.FailedUnauthorizedResponse
// @Failure 500 {object} jsonHttpResponse.FailedInternalServerErrorResponse
// @Security JWTToken
// @Router /api/v1/auth/change-password [post]
func (handler *authDelivery) changePassword(c *gin.Context) {
	var request authDTO.ReqUpdatePassword
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			jsonHttpResponse.NewFailedBadRequestResponse(c, validations)
			return
		}

		jsonHttpResponse.NewFailedBadRequestResponse(c, errBind.Error())
		return
	}

	userFromJWT, err := routeHelpers.GetUserFromJWTContext(c)
	if err != nil {
		if err == authmiddleware.ErrUserContextNotSet {
			jsonHttpResponse.InternalServerError(c, jsonHttpResponse.NewFailedResponse(err.Error()))
			return
		}
		jsonHttpResponse.Unauthorized(c, jsonHttpResponse.NewFailedResponse(err.Error()))
		return
	}

	//execute use case for change password
	err = handler.authUC.ChangePassword(userFromJWT.ID, request.PreviousPassword, request.NewPassword)
	if err != nil {
		if err == authUseCase.ErrPasswordMismatch ||
			err == authUseCase.ErrInvalidCredential {
			jsonHttpResponse.NewFailedUnauthorizedResponse(c, err.Error())
			return
		}

		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(
		c,
		authDTO.ResUpdatePassword{Mesage: "password succesfully changed"},
	)
	return
}

// Authentication - Reset password (Forgot Password) godoc
// @Summary Reset Password / Forget password
// @Description Reset user password by sending email to reset password form
// @ID authentication-reset-password-by-email
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param applicationData body authDTO.ReqUpdatePassword true "Chaneg password data"
// @Success 200 {object} jsonHttpResponse.SuccessResponse{data=authDTO.ResResetPasswordByEmail}
// @Failure 400 {object} jsonHttpResponse.FailedBadRequestResponse
// @Failure 401 {object} jsonHttpResponse.FailedUnauthorizedResponse
// @Failure 500 {object} jsonHttpResponse.FailedInternalServerErrorResponse
// @Security JWTToken
// @Router /api/v1/reset-password/send-mail [post]
func (handler *authDelivery) resetPasswordByEmailSend(c *gin.Context) {
	var request authDTO.ReqResetPasswordByEmail
	errBind := c.ShouldBind(&request)
	if errBind != nil {
		validations := requestvalidationerror.GetvalidationError(errBind)

		if len(validations) > 0 {
			jsonHttpResponse.NewFailedBadRequestResponse(c, validations)
			return
		}

		jsonHttpResponse.NewFailedBadRequestResponse(c, errBind.Error())
		return
	}

	_, err := handler.authUC.SendPasswordResetByEmail(request.Email)
	if err != nil {
		if err == authUseCase.ErrUserNotFound {
			jsonHttpResponse.NewFailedBadRequestResponse(c, err.Error())
			return
		}
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	res := authDTO.ResResetPasswordByEmail{
		Message: "email reset request successfully sent to " + request.Email,
	}
	jsonHttpResponse.NewSuccessfulOKResponse(c, res)
	return
}

// Verify reset password godoc
// @Summary Reset password verification
// @Description Verify reset password token
// @ID authentication-verify-reset-password
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Success 200 {object} jsonHttpResponse.SuccessResponse{data=authDTO.ResUpdatePassword}
// @Failure 400 {object} jsonHttpResponse.FailedBadRequestResponse
// @Failure 401 {object} jsonHttpResponse.FailedUnauthorizedResponse
// @Failure 500 {object} jsonHttpResponse.FailedInternalServerErrorResponse
// @Security JWTToken
// @Router /api/v1/auth/reset-password/verify-reset [post]
func (handler *authDelivery) verifyReset(c *gin.Context) {
	newPassword := c.PostForm("password")
	passwordConfirmation := c.PostForm("confirm-password")

	//password confirmation check
	if newPassword != passwordConfirmation {
		jsonHttpResponse.NewFailedBadRequestResponse(c, "password mismatch")
		return
	}

	userToken := c.PostForm("token")
	userTokenCollapse := strings.Split(userToken, "_")
	if len(userTokenCollapse) < 2 {
		jsonHttpResponse.NewFailedBadRequestResponse(c, "error invalid user token")
		return
	}

	userID, err := strconv.Atoi(userTokenCollapse[0])
	if err != nil {
		jsonHttpResponse.NewFailedBadRequestResponse(c, "error invalid user token")
		return
	}

	err = handler.authUC.ResetPasswordByToken(userID, userToken, newPassword)
	if err != nil {
		if err == authUseCase.ErrInvalidEmailVerificationToken {
			jsonHttpResponse.NewFailedUnauthorizedResponse(c, err.Error())
			return
		}
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	c.HTML(200, "auth/reset-password-success.html", nil)
	return
}
