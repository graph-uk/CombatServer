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

	combatServer := &server.CombatServer{}

	if config.MaxStoredSessions > 0 {
		sessionsRepo.DeleteOldSessions(config.MaxStoredSessions)
	}

	err = combatServer.Start(config)
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
