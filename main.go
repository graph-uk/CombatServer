package main

import (
	"net/http"

	"github.com/graph-uk/combat-server/server"
)

func main() {
	// serve static files for HTML pages (jquery, css, images, etc)
	// to update static files - run packBinData.cmd

	http.Handle("/bindata/", http.StripPrefix("/bindata/", http.FileServer(assetFS())))
	combatServer := server.NewCombatServer()
	combatServer.Serve()
}
