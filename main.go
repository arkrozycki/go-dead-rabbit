package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var LOG_TYPE = os.Getenv("LOG_TYPE")
var ENVIRONMENT = os.Getenv("ENVIRONMENT")

// init
func init() {
	// configure the zero logger
	LOG_LEVEL, _ := zerolog.ParseLevel(LOG_TYPE)
	fmt.Println(LOG_TYPE)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(LOG_LEVEL)

	// while in dev env simply output to console
	if ENVIRONMENT == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// main
func main() {
	log.Debug().Msg("Go Dead Rabbit Started")
}
