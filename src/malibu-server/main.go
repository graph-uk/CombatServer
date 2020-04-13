package main

import (
	"fmt"
	"os"

	"malibu-server/server"
	"malibu-server/utils"
)

func main() {
	config := utils.GetApplicationConfig()

	db := utils.GetDB()
	defer db.Close()

	malibuServer := &server.MalibuServer{}
	err := malibuServer.Start(config, db)
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
