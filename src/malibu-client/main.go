package main

import (
	"fmt"
	"os"

	"malibu-client/client"
)

func main() {

	client, err := malibuClient.NewMalibuClient()
	if err != nil {
		fmt.Println("Cannot init malibu client")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sessionID, err := client.CreateNewSession()
	if err != nil {
		fmt.Println("Cannot create malibu session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(client.GetSessionResult(sessionID))
}
