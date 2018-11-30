package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/combat-server/utils"

	"github.com/graph-uk/combat-server/data/repositories"

	"github.com/graph-uk/combat-server/server"
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
