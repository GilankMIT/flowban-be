package userModuleInit

import (
	"flowban/helper/sendgridMailHelper"
	"flowban/middleware/authmiddleware"
	"flowban/module/user/userDelivery"
	"flowban/module/user/userRepository"
	"flowban/module/user/userUsecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(r *gin.Engine, db *gorm.DB, authMiddleware authmiddleware.AuthMiddleware,
	mailService sendgridMailHelper.SendGridMailService) {
	userRepo := userRepository.NewUserRepository(db)
	userUseCase := userUsecase.NewUserRepositoryUseCase(userRepo, mailService)
	userDelivery.NewUserHTTPDelivery(r, userUseCase, authMiddleware)
}
