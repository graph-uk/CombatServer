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
	repo.Apply()

	config := utils.GetApplicationConfig()

	malibuServer := &server.MalibuServer{}

	err := malibuServer.Start(config)
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
