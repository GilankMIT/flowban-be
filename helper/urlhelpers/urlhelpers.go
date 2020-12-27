package urlhelpers

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	BaseURL string
)

func init() {
	err := godotenv.Load("conf/application.env")
	if err != nil {
		log.Error().Msg("Failed read configuration database")
		return
	}

	BaseURL = os.Getenv("base_url")
}

func GetBaseURL() string {
	return BaseURL
}
