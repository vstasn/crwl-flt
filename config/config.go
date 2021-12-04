package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var AppConfig = getConfig()

type Config struct {
	FltURL           string `envconfig:"FLT_URL" required:"true"`
	PgURL            string `envconfig:"DATABASE_URL" required:"true"`
	TelegramApiToken string `envconfig:"TELEGRAM_API_TOKEN" required:"true"`
	TelegramChatId   int64  `envconfig:"TELEGRAM_CHAT_ID" required:"true"`
}

func getConfig() Config {
	var cnf Config
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Fatal(err)
	}
	return cnf
}
