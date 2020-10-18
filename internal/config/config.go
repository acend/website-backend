package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

var Values struct {
	MailchimpAPIKey string `required:"true" split_words:"true"`
	SMTP            struct {
		Host     string `required:"true"`
		Port     int    `required:"true"`
		Username string `required:"true"`
		Password string `required:"true"`
		From     string `required:"true"`
		To       string `required:"true"`
	}
}

func Configure() {
	err := envconfig.Process("", &Values)
	if err != nil {
		log.Fatal(err)
	}
}
