package worker

import (
	"training/hello-cadence/app/adapters/cadenceAdapter"

	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func StartWorkers(adapter *cadenceAdapter.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: adapter.Scope,
		Logger:       adapter.Logger,
	}

	cadenceWorker := worker.New(
		adapter.ServiceClient, adapter.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		adapter.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}
