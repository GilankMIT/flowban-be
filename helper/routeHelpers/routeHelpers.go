package routeHelpers

import (
	"errors"
	"flowban/middleware/authmiddleware"
	"flowban/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var (
	ErrContextNotExist                 = errors.New("user context not exist")
	ErrParsingUserModel                = errors.New("error parsing user model")
	ErrInvalidPageHeaderFormat         = errors.New("invalid page header format")
	ErrInvalidPerPageLimitHeaderFormat = errors.New("invalid per page header format")
)

func GetUserFromJWTContext(c *gin.Context) (*model.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, ErrContextNotExist
	}
	if user == authmiddleware.ErrUserContextNotSet {
		return nil, ErrContextNotExist
	}

	userModel, ok := user.(*model.User)
	if !ok {
		return nil, ErrParsingUserModel
	}

	return userModel, nil
}

func GetPaginationInfo(c *gin.Context) (perPageLimit, page int, err error) {
	//get pagination header

	perPageLimitHeader := c.GetHeader("x-per-page-limit")
	if perPageLimitHeader != "" {
		perPageLimit, err = strconv.Atoi(perPageLimitHeader)
		if err != nil {
			return 0, 0, ErrInvalidPerPageLimitHeaderFormat
		}
	}

	pageHeader := c.GetHeader("x-page")
	if pageHeader != "" {
		page, err = strconv.Atoi(pageHeader)
		if err != nil {
			return 0, 0, ErrInvalidPageHeaderFormat
		}
	}

	return perPageLimit, page, nil
}

func GetOrderFieldInfo(c *gin.Context) (orderField string, isExist bool) {
	orderField = c.GetHeader("x-order-field")

	if orderField != "" {
		return orderField, true
	}

	return orderField, false
}

func GetFilterColumn(c *gin.Context) (filterColumn string) {
	return c.GetHeader("x-filter-column")
}

func GetFilterColumnValue(c *gin.Context) (filterColumn string) {
	return c.GetHeader("x-filter-column-value")
}

func GetDateEarliest(c *gin.Context) (dateEarliestVal *time.Time, err error) {
	dateEarliestHeader := c.GetHeader("x-date-earliest")
	if dateEarliestHeader != "" {
		dateEarliest, err := time.Parse("2006-01-02", dateEarliestHeader)
		return &dateEarliest, err
	}
	return nil, nil
}

func GetDateLatest(c *gin.Context) (dateLatestVal *time.Time, err error) {
	dateLatestHeader := c.GetHeader("x-date-latest")
	if dateLatestHeader != "" {
		dateLatest, err := time.Parse("2006-01-02", dateLatestHeader)
		return &dateLatest, err
	}

	return nil, nil
}
