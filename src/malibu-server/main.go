package main

import (
	"fmt"
	"os"

	"malibu-server/utils"

	"malibu-server/data/repositories"

	"malibu-server/server"
)

func main() {
	repo := &repositories.Migrations{}
	sessionsRepo := &repositories.Sessions{}

	config := utils.GetApplicationConfig()

	err := repo.Apply()

	if err != nil {
		panic(err)
	}

	malibuServer := &server.MalibuServer{}

	if config.MaxStoredSessions > 0 {
		sessionsRepo.DeleteOldSessions(config.MaxStoredSessions)
	}

	err = malibuServer.Start(config)
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
