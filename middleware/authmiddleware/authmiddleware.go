package authmiddleware

import "github.com/gin-gonic/gin"

type AuthMiddleware interface {
	AuthorizeJWT(roles string) gin.HandlerFunc
	AuthorizeJWTWithUserContext(roles string) gin.HandlerFunc
}
