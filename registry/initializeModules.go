package registry

import (
	"flowban/middleware/authmiddleware"
	"flowban/module/user/userRepository"
	"flowban/moduleInit/authModuleInit"
	"flowban/moduleInit/scrumProjectModuleInit"
	"flowban/moduleInit/userModuleInit"
)

var (
	authMiddleware authmiddleware.AuthMiddleware
)

//initializeDomainModules calls the domain module routes in folder moduleInit/*
func (reg *AppRegistry) initializeDomainModules() {
	reg.initializeMiddleware()
	authModuleInit.InitModule(reg.httpHandler.GetRouteEngine(), reg.dbConn, reg.emailService, authMiddleware)
	userModuleInit.InitModule(reg.httpHandler.GetRouteEngine(), reg.dbConn, authMiddleware, reg.emailService)
	scrumProjectModuleInit.InitModule(reg.httpHandler.GetRouteEngine(), reg.dbConn, authMiddleware)
}

//initialize middleware is used to init middlewares and
//inject the needed dependency
func (reg *AppRegistry) initializeMiddleware() {
	userRepo := userRepository.NewUserRepository(reg.dbConn)
	authMiddleware = authmiddleware.NewAuthMiddleware(userRepo)
}
