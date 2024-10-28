package server

import (
	"finance-service/config"
	"log"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

}
