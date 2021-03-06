package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var LOG_TYPE = os.Getenv("LOG_TYPE")
var ENVIRONMENT = os.Getenv("ENVIRONMENT")
var CONFIG_FILE = "config.yml"

type Config struct {
	Connection   ConnectionConfig
	Listener     ListenerConfig
	Publisher    PublisherConfig
	Notification NotificationConfig
	Datastore    DatastoreConfig
}

type ConnectionConfig struct {
	Server   string
	Port     string
	Vhost    string
	User     string
	Password string
}

type ListenerConfig struct {
	Queue QueueConfig
}

type QueueConfig struct {
	Name string
}

type PublisherConfig struct {
	Exchange ExchangeConfig
}

type ExchangeConfig struct {
	Name string
}

type NotificationConfig struct {
	Mailgun MailgunConfig
}

type MailgunConfig struct {
	ApiKey  string
	BaseUrl string
	Domain  string
	From    string
	To      string
}

type DatastoreConfig struct {
	Mongodb MongodbConfig
}

type MongodbConfig struct {
	Database   string
	Collection string
	Uri        string
}

var Conf Config

// init
// Sets up the configuration
// reads yml and env variables into object
func init() {
	// configure the zero logger
	LOG_LEVEL, _ := zerolog.ParseLevel(LOG_TYPE)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(LOG_LEVEL)

	// while in dev env simply output to console
	if ENVIRONMENT == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	viper.SetConfigFile(CONFIG_FILE)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Err(err).Msg("Error reading config file")
	}

	// unmarshall into viper configuration object
	err := viper.Unmarshal(&Conf)
	if err != nil {
		log.Err(err).Msg("Unable to decode into struct")
	}

	// collect environment variables
	Conf.Connection.Server = viper.GetString("RABBIT_SERVER")
	Conf.Connection.Port = viper.GetString("RABBIT_PORT")
	Conf.Connection.Vhost = viper.GetString("RABBIT_VHOST")
	Conf.Connection.User = viper.GetString("RABBIT_USER")
	Conf.Connection.Password = viper.GetString("RABBIT_PASSWORD")

	// set up additional env variables for mailgun
	Conf.Notification.Mailgun.ApiKey = viper.GetString("MAILGUN_API_KEY")
	Conf.Notification.Mailgun.Domain = viper.GetString("MAILGUN_API_DOMAIN")

	RABBIT_STRING := fmt.Sprintf("%s@%s:%s/%s", Conf.Connection.User, viper.GetString("RABBIT_SERVER"), viper.GetString("RABBIT_PORT"), viper.GetString("RABBIT_VHOST"))
	MAILGUN_STRING := fmt.Sprintf("%s@%s/%s", Conf.Notification.Mailgun.ApiKey, Conf.Notification.Mailgun.BaseUrl, Conf.Notification.Mailgun.Domain)
	log.Info().
		Str("RABBIT", RABBIT_STRING).
		Str("MAILGUN", MAILGUN_STRING).
		Str("LISTENER QUEUE", Conf.Listener.Queue.Name).
		Str("PUBLISHER EXCHANGE", Conf.Publisher.Exchange.Name).
		Msg("CONFIG")
}
