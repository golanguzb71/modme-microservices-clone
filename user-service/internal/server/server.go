package server

import (
	"fmt"
	"log"
	"user-service/config"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Print(cfg)

}
