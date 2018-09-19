package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/combat-server/data/repositories"

	"github.com/graph-uk/combat-server/server"
)

func main() {
	repo := &repositories.Migrations{}
	repo.Apply()

	combatServer, err := server.NewCombatServer()
	if err != nil {
		fmt.Println("Cannot init combat server")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	combatServer.DeleteOldSessions()

	err = combatServer.Start()
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
