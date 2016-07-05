package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/graph-uk/combat-server/server"
)

func main() {
	// serve static files for HTML pages (jquery, css, images, etc)
	// to update static files - run packBinData.cmd
	http.Handle("/bindata/", http.StripPrefix("/bindata/", http.FileServer(assetFS())))

	combatServer, err := server.NewCombatServer()
	if err != nil {
		fmt.Println("Cannot init combat server")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	combatServer.DeleteOldSessions()
	//os.Exit(1)

	err = combatServer.Serve()
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
