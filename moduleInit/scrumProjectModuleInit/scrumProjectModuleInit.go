package scrumProjectModuleInit

import (
	"flowban/middleware/authmiddleware"
	"flowban/module/scrumProject/scrumProjectDelivery"
	"flowban/module/scrumProject/scrumProjectRepository"
	"flowban/module/scrumProject/scrumProjectUseCase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(r *gin.Engine, db *gorm.DB, middleware authmiddleware.AuthMiddleware) {
	scrumProjRepo := scrumProjectRepository.NewScrumProjectRepository(db)
	scrumProjUC := scrumProjectUseCase.NewScrumProjectRepositoryUseCase(scrumProjRepo)
	scrumProjectDelivery.NewHTTPDelivery(r, scrumProjUC, middleware)
}
