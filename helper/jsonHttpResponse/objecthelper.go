package jsonHttpResponse

import (
	"github.com/gin-gonic/gin"
)

//type helpers. Used as an example value in swagger documentation

type SuccessResponse struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

type FailedResponse struct {
	Status  string      `json:"status" example:"failed"`
	Message interface{} `json:"message"`
}

//responses helper
type FailedUnauthorizedResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"invalid user token"`
}

type FailedBadRequestResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"bad request"`
}

type FailedUnprocessableEntityResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"unprocessable entity"`
}

type FailedInternalServerErrorResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"internal server error"`
}

type FailedNotFoundResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"data not found"`
}

func NewSuccessResponse(payload interface{}) SuccessResponse {
	return SuccessResponse{Status: "success", Data: payload}
}

func NewFailedResponse(message interface{}) FailedResponse {
	return FailedResponse{Status: "failed", Message: message}
}

func NewSuccessfulOKResponse(c *gin.Context, payload interface{}) {
	OK(c, NewSuccessResponse(payload))
}

func NewFailedBadRequestResponse(c *gin.Context, message interface{}) {
	BadRequest(c, NewFailedResponse(message))
}

func NewNoContentResponse(c *gin.Context, message interface{}) {
	NoContent(c, NewFailedResponse(message))
}
func NewNotFoundResponse(c *gin.Context, message interface{}) {
	NotFound(c, NewFailedResponse(message))
}

func NewFailedUnauthorizedResponse(c *gin.Context, message interface{}) {
	Unauthorized(c, NewFailedResponse(message))
}

func NewFailedUnprocessableResponse(c *gin.Context, message interface{}) {
	Unprocessable(c, NewFailedResponse(message))
}

func NewFailedForbiddenResponse(c *gin.Context, message interface{}) {
	Forbidden(c, NewFailedResponse(message))
}

func NewFailedInternalServerResponse(c *gin.Context, message interface{}) {
	InternalServerError(c, NewFailedResponse(message))
}

func NewFailedConflictResponse(c *gin.Context, message interface{}) {
	Conflict(c, NewFailedResponse(message))
}
