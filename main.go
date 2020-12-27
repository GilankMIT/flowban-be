package main

import (
	"flowban/registry"
	"github.com/swaggo/gin-swagger/example/basic/docs"
)

func main() {
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	appRegistry := registry.NewAppRegistry()
	appRegistry.StartServer()
}
