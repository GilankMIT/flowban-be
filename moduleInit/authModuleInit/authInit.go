package authModuleInit

import (
	"flowban/helper/sendgridMailHelper"
	"flowban/middleware/authmiddleware"
	"flowban/module/auth/authDelivery"
	"flowban/module/auth/authUseCase"
	"flowban/module/user/userRepository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(r *gin.Engine, db *gorm.DB, emailService sendgridMailHelper.SendGridMailService,
	authMiddleware authmiddleware.AuthMiddleware) {
	userRepo := userRepository.NewUserRepository(db)
	authUseCase := authUseCase.NewAuthUseCase(userRepo, emailService)
	authDelivery.NewAuthHTTPDelivery(r, authUseCase, authMiddleware)
}
