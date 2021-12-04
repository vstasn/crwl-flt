package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var AppConfig = getConfig()

type Config struct {
	FltURL           string `envconfig:"FLT_URL" required:"true"`
	PgURL            string `envconfig:"POSTGRES_URL" required:"true"`
	TelegramApiToken string `envconfig:"TELEGRAM_API_TOKEN" required:"true"`
	TelegramChatId   int64  `envconfig:"TELEGRAM_API_TOKEN" required:"true"`
}

func getConfig() Config {
	var cnf Config
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Fatal(err)
	}
	return cnf
}
