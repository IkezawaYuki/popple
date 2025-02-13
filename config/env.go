package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Environment struct {
	RedisAddr              string `envconfig:"REDIS_ADDR"`
	DatabaseUser           string `envconfig:"DATABASE_USER"`
	DatabasePort           string `envconfig:"DATABASE_PORT"`
	DatabasePass           string `envconfig:"DATABASE_PASSWORD"`
	DatabaseName           string `envconfig:"DATABASE_SCHEME"`
	DatabaseHost           string `envconfig:"DATABASE_HOST"`
	AccessSecretKey        string `envconfig:"ACCESS_SECRET_KEY"`
	WordpressAdminEmail    string `envconfig:"WORDPRESS_ADMIN_EMAIL"`
	WordpressAdminPassword string `envconfig:"WORDPRESS_ADMIN_PASSWORD"`
	GraphApiURL            string `envconfig:"GRAPH_API_URL"`
	SlackWebhookURL        string `envconfig:"SLACK_WEBHOOK_URL"`
	RodutKey               string `envconfig:"RODUT_KEY"`
}

var Env Environment

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := envconfig.Process("", &Env); err != nil {
		log.Fatal(err)
	}
}
