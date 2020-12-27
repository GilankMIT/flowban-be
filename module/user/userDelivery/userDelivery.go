package userDelivery

import (
	"flowban/helper/jsonHttpResponse"
	"flowban/helper/requestvalidationerror"
	"flowban/helper/routeHelpers"
	"flowban/middleware/authmiddleware"
	"flowban/module/user"
	"flowban/module/user/userDTO"
	"flowban/module/user/userUsecase"
	"github.com/gin-gonic/gin"
	"strconv"
)

type userDelivery struct {
	userUC user.UserUseCase
}

func NewUserHTTPDelivery(r *gin.Engine, userUC user.UserUseCase,
	authMiddleware authmiddleware.AuthMiddleware) {

	handler := userDelivery{
		userUC: userUC,
	}

	authorizedAdmin := r.Group("api/v1/user",
		authMiddleware.AuthorizeJWTWithUserContext("admin"))
	{
		authorizedAdmin.GET("", handler.getAll)
		authorizedAdmin.POST("", handler.addNew)
		authorizedAdmin.DELETE("", handler.removeByID)
		authorizedAdmin.GET("roles", handler.getAllRoles)
		authorizedAdmin.GET("by-id/:id", handler.getUserByID)
	}

	authorizedUser := r.Group("api/v1/user",
		authMiddleware.AuthorizeJWTWithUserContext("admin|user"))
	{
		authorizedUser.GET("me", handler.getMyUserData)
		authorizedUser.POST("me/update", handler.updateData)
	}
}

func (handler *userDelivery) getAllRoles(c *gin.Context) {
	roles, err := handler.userUC.GetAllRole()
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, roles)
	return
}

func (handler *userDelivery) getMyUserData(c *gin.Context) {
	userPrincipal, err := routeHelpers.GetUserFromJWTContext(c)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}
	roles, err := handler.userUC.GetByID(userPrincipal.ID)
	if err != nil {
		if err == userUsecase.ErrUserNotFound {
			jsonHttpResponse.NewNotFoundResponse(c, err.Error())
			return
		}
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, roles)
	return
}

func (handler *userDelivery) getUserByID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)

	userData, err := handler.userUC.GetByID(userID)
	if err != nil {
		if err == userUsecase.ErrUserNotFound {
			jsonHttpResponse.NewNotFoundResponse(c, err.Error())
			return
		}
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, userData)
	return
}

// Update users godoc
// @Summary Update own user
// @ID update-own-user
// @Tags User
// @Accept  json
// @Produce  json
// @Param UserData body userDTO.ReqUpdateOwnUser true "Usr Data"
// @Success 200 {object} jsonHttpResponse.SuccessResponse{data=[]model.User}
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/user/me/update [post]
func (handler *userDelivery) updateData(c *gin.Context) {
	var request userDTO.ReqUpdateOwnUser
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

	userData, err := routeHelpers.GetUserFromJWTContext(c)
	if err != nil {
		jsonHttpResponse.Unauthorized(c, err.Error())
		return
	}

	_, err = handler.userUC.UpdateCredentialData(userData.ID, request.Password,
		request.FirstName, request.LastName)
	if err != nil {
		if err == userUsecase.ErrUserNotFound {
			jsonHttpResponse.NewFailedBadRequestResponse(c, err.Error())
			return
		}

		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	res := userDTO.ResRemoveUserByID{Message: "user data successfully modified"}
	jsonHttpResponse.NewSuccessfulOKResponse(c, res)
	return
}

// Get All Users godoc
// @Summary Get All Users
// @ID user-get-all
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} jsonHttpResponse.SuccessResponse{data=[]model.User}
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/user [get]
func (handler *userDelivery) getAll(c *gin.Context) {
	userData, err := handler.userUC.GetAll()
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, userData)
	return
}

// Add New User godoc
// @Summary Add New User
// @ID user-add
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} model.User
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/user/add [post]
func (handler *userDelivery) addNew(c *gin.Context) {
	var request userDTO.ReqAddNewUser
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

	userData, err := handler.userUC.AddNew(request)
	if err != nil {
		if err == userUsecase.ErrEmailAlreadyExist {
			jsonHttpResponse.NewFailedConflictResponse(c, err.Error())
			return
		}
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, userData)
	return
}

// Remove User By ID godoc
// @Summary Get User by ID
// @ID user-get-by-id
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} userDTO.ResRemoveUserByID
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /api/v1/user [delete]
func (handler *userDelivery) removeByID(c *gin.Context) {
	var request userDTO.ReqRemoveUserByID
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

	err := handler.userUC.RemoveByID(request.UserID)
	if err != nil {
		if err == userUsecase.ErrUserNotFound {
			jsonHttpResponse.NewFailedBadRequestResponse(c, err.Error())
			return
		}
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	res := userDTO.ResRemoveUserByID{Message: "user successfully deleted"}
	jsonHttpResponse.NewSuccessfulOKResponse(c, res)
	return
}
