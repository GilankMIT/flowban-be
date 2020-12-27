package scrumProjectDelivery

import (
	"flowban/helper/jsonHttpResponse"
	"flowban/helper/requestvalidationerror"
	"flowban/helper/routeHelpers"
	"flowban/middleware/authmiddleware"
	"flowban/model"
	"flowban/module/scrumProject"
	"flowban/module/scrumProject/scrumProjectDTO"
	"github.com/gin-gonic/gin"
	"strconv"
)

type httpDelivery struct {
	scrumProjectUC scrumProject.ScrumProjectUseCase
}

func NewHTTPDelivery(r *gin.Engine, scrumProjectUC scrumProject.ScrumProjectUseCase, authMiddleware authmiddleware.AuthMiddleware) {
	handler := httpDelivery{scrumProjectUC: scrumProjectUC}

	authorizedUser := r.Group("api/v1/scrum-project", authMiddleware.AuthorizeJWTWithUserContext("user|admin"))
	{
		authorizedUser.GET("my-project", handler.getProjectByJWT)
		authorizedUser.GET("by-id/:project_id", handler.getProjectByID)
		authorizedUser.GET("active-issues/:project_id", handler.getActiveIssuesByProjectID)
		authorizedUser.GET("boards/:project_id", handler.getBoardByProjectID)
		authorizedUser.GET("issues/by-id/:issue-id", handler.getIssueDetail)

		authorizedUser.POST("create-project", handler.createNewProject)
		authorizedUser.POST("create-sprint", handler.createSprint)
		authorizedUser.POST("create-issue", handler.createIssue)
		authorizedUser.POST("create-board", handler.createBoard)
		authorizedUser.POST("move-issue-from-backlog", handler.moveIssueFromBacklog)
		authorizedUser.POST("move-issue-to-backlog", handler.moveIssueToBacklog)
		authorizedUser.POST("move-issue", handler.moveIssue)

		authorizedUser.DELETE("delete-issue", handler.deleteIssue)

	}
}

func (handler *httpDelivery) moveIssue(c *gin.Context) {
	var request scrumProjectDTO.ReqMoveIssue
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

	issue, err := handler.scrumProjectUC.MoveIssue(request.IssueID, request.BoardID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, issue)
	return
}

func (handler *httpDelivery) getProjectByID(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	project, err := handler.scrumProjectUC.GetByID(projectID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, project)
	return
}

func (handler *httpDelivery) createNewProject(c *gin.Context) {
	var request scrumProjectDTO.ReqAddNewProject
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

	//build model
	userPrinc, err := routeHelpers.GetUserFromJWTContext(c)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	scrumProjectData := model.ScrumProject{
		Name:        request.Name,
		Description: request.Description,
		UserID:      userPrinc.ID,
		ImageURL:    request.ImageURL,
		Acronym:     request.Acronym,
	}

	data, err := handler.scrumProjectUC.AddNew(scrumProjectData)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, data)
	return
}

func (handler *httpDelivery) getProjectByJWT(c *gin.Context) {
	//get user
	userPrincipal, err := routeHelpers.GetUserFromJWTContext(c)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	projects, err := handler.scrumProjectUC.GetByUserID(userPrincipal.ID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, projects)
	return
}

func (handler *httpDelivery) getActiveIssuesByProjectID(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	issues, err := handler.scrumProjectUC.GetActiveIssueByProjectID(projectID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, issues)
	return
}

func (handler *httpDelivery) getBoardByProjectID(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
	}

	boards, err := handler.scrumProjectUC.GetProjectBoards(projectID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, boards)
	return
}

func (handler *httpDelivery) getIssueDetail(c *gin.Context) {

}

func (handler *httpDelivery) createSprint(c *gin.Context) {
	var request scrumProjectDTO.ReqCreateNewSprint
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

	sprintData, err := handler.scrumProjectUC.AddNewSprint(request)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, sprintData)
	return
}

func (handler *httpDelivery) createIssue(c *gin.Context) {
	var request scrumProjectDTO.ReqCreateNewIssue
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

	newIssue, err := handler.scrumProjectUC.AddNewIssue(request)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, newIssue)
	return
}

func (handler *httpDelivery) createBoard(c *gin.Context) {

}

func (handler *httpDelivery) moveIssueFromBacklog(c *gin.Context) {
	var request scrumProjectDTO.ReqMoveIssueInBacklog
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

	savedIssue, err := handler.scrumProjectUC.MoveIssueFromBacklog(request.IssueID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, savedIssue)
	return
}

func (handler *httpDelivery) moveIssueToBacklog(c *gin.Context) {
	var request scrumProjectDTO.ReqMoveIssueInBacklog
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

	savedIssue, err := handler.scrumProjectUC.MoveIssueToBacklog(request.IssueID)
	if err != nil {
		jsonHttpResponse.NewFailedInternalServerResponse(c, err.Error())
		return
	}

	jsonHttpResponse.NewSuccessfulOKResponse(c, savedIssue)
	return
}

func (handler *httpDelivery) deleteIssue(c *gin.Context) {

}
