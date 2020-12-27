package registry

import (
	"flowban/helper/dsnbuilder"
	"flowban/helper/sendgridMailHelper"
	"flowban/registry/serverhandler"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const (
	EnvPath = "conf/.env"
)

var (
	// Application configuration variable
	//appVersion string
	appPort string

	// DB configuration variable
	dbHost     string
	dbPort     string
	dbDatabase string
	dbUsername string
	dbPassword string
	dbDialect  string
	dbDsn      string

	oauthURL      string
	oauthUsername string
	oauthPassword string
)

type AppRegistry struct {
	dbConn *gorm.DB

	httpHandler  *serverhandler.HttpHandler
	emailService sendgridMailHelper.SendGridMailService
}

//NewAppRegistry will return new object for App Registry
func NewAppRegistry() *AppRegistry {
	return &AppRegistry{}
}

//StartServer will do the server initialization
func (reg *AppRegistry) StartServer() {
	reg.initializeAppRegistry()

	//Run Swagger
	log.Info().Msg("Swagger run on /docs/swagger/index.html")
	reg.httpHandler.RunSwaggerMiddleware()

	//Run HTTP Server
	appVersion := "0.1.9"
	log.Info().Msg("Last Update : " + time.Now().Format("2006-01-02 15:04:05"))
	log.Info().Msg("REST API Service Running version " + appVersion + " at port : " + appPort)
	if errHTTP := reg.httpHandler.RunHttpServer(); errHTTP != nil {
		log.Error().Msg(errHTTP.Error())
	}

	//Close connection
	defer func() {
		db, err := reg.dbConn.DB()
		if err != nil {
			log.Error().Msg(err.Error())
		}

		err = db.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()

}

func (registry *AppRegistry) initializeAppRegistry() {
	//Initialize Logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stdout,
			NoColor: false,
		},
	)

	printAsciiArt()

	err := initializeEnv()
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}

	err = registry.initializeDatabase()
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}

	err = registry.initializeDependency()
	if err != nil {
		log.Error().Msg(err.Error())
		panic("app exit on error")
	}

	registry.initializeHandler()

	//Initialize modules from folder modules,
	//injecting dependencies and running preloading script if exist
	registry.initializeDomainModules()
}

func initializeEnv() error {
	err := godotenv.Load(EnvPath)
	if err != nil {
		log.Error().Msg("Failed to read configuration database")
		return err
	}

	if os.Getenv("project_env") == "DEV" || os.Getenv("project_env") == "" { //project_env is Development

		//appVersion = os.Getenv("application.version")
		appPort = os.Getenv("application.port")

		//specify db dialect
		dbDialect = os.Getenv("db.dialect")

		dbHost = os.Getenv("db." + dbDialect + ".host")
		dbPort = os.Getenv("db." + dbDialect + ".port")
		dbUsername = os.Getenv("db." + dbDialect + ".username")
		dbPassword = os.Getenv("db." + dbDialect + ".password")
		dbDatabase = os.Getenv("db." + dbDialect + ".dbname")

		//initialize Oauth env
		oauthURL = os.Getenv("oauth.url")
		oauthUsername = os.Getenv("oauth.username")
		oauthPassword = os.Getenv("oauth.password")
	}
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic(err)
	}
	dbDsn, err = dsnbuilder.New(dbHost, dbPortInt, dbUsername, dbPassword, dbDatabase).Build(dbDialect)
	if err != nil {
		panic(err)
	}

	return nil
}

func (reg *AppRegistry) initializeHandler() {
	//Register HTTP Server Handler
	reg.httpHandler = serverhandler.NewHTTPHandler(":" + appPort)
}

func (reg *AppRegistry) initializeDependency() error {
	sendGridAPIKey := os.Getenv("sendgrid.apikey")
	reg.emailService = sendgridMailHelper.NewEmailService(sendGridAPIKey)
	return nil
}
