package main

import (
	"fmt"
	"os"

	"malibu-worker/worker"
)

func main() {
	worker, err := worker.NewCombatWorker()
	if err != nil {
		fmt.Println("Cannot init combat worker")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		worker.Process()
	}
}
