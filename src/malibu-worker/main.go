package main

import (
	"fmt"
	"os"

	"malibu-worker/worker"
)

func main() {
	worker, err := worker.NewMalibuWorker()
	if err != nil {
		fmt.Println("Cannot init malibu worker")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		worker.Process()
	}
}
