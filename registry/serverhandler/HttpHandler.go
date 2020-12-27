package serverhandler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
	"regexp"
	"time"
)

type HttpHandler struct {
	r       *gin.Engine
	address string
}

func NewHTTPHandler(adress string) *HttpHandler {
	r := gin.New()
	var rxURL = regexp.MustCompile(`^/regexp\d*`)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST", "DELETE", "GET", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Origin", "Content-Type", "content-type", "Authorization", "userid",
			"fromDate", "toDate", "x-page", "x-per-page-limit", "x-order-field"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           720 * time.Hour,
	}))

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.=
	subLog := zerolog.New(os.Stdout).With().Logger()

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
	r.Use(logger.SetLogger())

	r.Use(logger.SetLogger(logger.Config{
		Logger:         &subLog,
		UTC:            true,
		SkipPath:       []string{"/skip"},
		SkipPathRegexp: rxURL,
	}))

	r.Use(gin.Recovery())

	//Maximum memory limit for Multipart forms
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	//r.Static("/assets/", "./assets/")

	return &HttpHandler{r, adress}
}

func (http *HttpHandler) GetRouteEngine() *gin.Engine {
	return http.r
}

func (http *HttpHandler) RunSwaggerMiddleware() {
	http.r.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (http *HttpHandler) RunHttpServer() error {
	if errHTTP := http.r.Run(http.address); errHTTP != nil {
		return errHTTP
	}
	return nil
}
