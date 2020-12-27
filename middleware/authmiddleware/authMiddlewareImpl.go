package authmiddleware

import (
	"errors"
	"flowban/helper/jsonHttpResponse"
	"flowban/helper/jwthelper"
	"flowban/module/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrInvalidToken      = errors.New("authmiddleware.authMiddlewareImpl : invalid token")
	ErrUserContextNotSet = errors.New("authmiddleware.authMiddlewareImpl : user context not set")
)

type authMiddleware struct {
	userRepo user.UserRepository
}

func NewAuthMiddleware(userRepo user.UserRepository) AuthMiddleware {
	return &authMiddleware{userRepo: userRepo}
}

//AuthorizeJWTWithUserContext - Authorize JWT with User Context (Need to look up for user in DB in every request)
func (auth *authMiddleware) AuthorizeJWTWithUserContext(roles string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		jwtClaims, err := extractAndValidateJWT(bearerToken, roles)
		if err != nil {
			jsonHttpResponse.NewFailedForbiddenResponse(c, err.Error())
			c.Abort()
			return
		}

		//get user by id
		userData, err := auth.userRepo.GetByID(jwtClaims.Id, true)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				jsonHttpResponse.NewFailedUnauthorizedResponse(c, "authMiddleware.authMiddlewareImpl "+err.Error())
				c.Abort()
				return
			}

			jsonHttpResponse.NewFailedInternalServerResponse(c, "authMiddleware.authMiddlewareImpl "+err.Error())
			c.Abort()
			return
		}

		c.Set("user", userData)
		c.Next()
		return
	}
}

//AuthorizeJWT - Authorize JWT without User Context (No need to look up for user in DB in every request)
func (auth *authMiddleware) AuthorizeJWT(roles string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		_, err := extractAndValidateJWT(bearerToken, roles)
		if err != nil {
			jsonHttpResponse.NewFailedForbiddenResponse(c, err.Error())
			c.Abort()
			return
		}

		c.Set("user", ErrUserContextNotSet)
		c.Next()
		return
	}
}

func extractAndValidateJWT(bearerToken, roles string) (*jwthelper.CustomClaims, error) {
	if bearerToken == "" {
		return nil, ErrInvalidToken
	}

	//Extract JWT Token from Bearer
	jwtTokenSplit := strings.Split(bearerToken, "Bearer ")
	if jwtTokenSplit[1] == "" {
		return nil, ErrInvalidToken
	}
	jwtToken := jwtTokenSplit[1]

	jwtTokenClaims, err := jwthelper.VerifyTokenWithClaims(jwtToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	jwtRolesArray := strings.Split(jwtTokenClaims.Role, "|")

	rolesArray := strings.Split(roles, "|")

	for _, role := range rolesArray {
		for _, jwtRole := range jwtRolesArray {
			if jwtRole == role {
				return jwtTokenClaims, nil
			}
		}
	}

	return nil, ErrInvalidToken
}
