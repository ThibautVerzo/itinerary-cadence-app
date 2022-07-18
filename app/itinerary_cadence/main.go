package main

import (
	"fmt"

	"training/hello-cadence/app/adapters/cadenceAdapter"
	"training/hello-cadence/app/config"
	"training/hello-cadence/app/itinerary_cadence/worker"
	"training/hello-cadence/app/itinerary_cadence/workflows"
)

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	worker.StartWorkers(&cadenceClient, workflows.TaskListName)
	// The workers are supposed to be long running process that should not exit.
	select {}
}
