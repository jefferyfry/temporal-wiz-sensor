package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/jefferyfry/temporal-wiz-sensor/activity"
	"github.com/jefferyfry/temporal-wiz-sensor/workflow"
)

func main() {
	temporalCl, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer temporalCl.Close()

	// Register the activity with the Temporal client
	sensorWorker := worker.New(temporalCl, "wiz-sensor-task-queue", worker.Options{})
	sensorWorker.RegisterWorkflow(workflow.DeploySensorWorkflow)
	sensorWorker.RegisterActivity(activity.InstallWizSensorActivity)

	log.Printf("Starting worker for task queue 'wiz-sensor-task-queue'...")
	errWorker := sensorWorker.Run(nil)
	if errWorker != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
