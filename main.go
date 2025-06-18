package main

import (
	"context"
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

	errWorker := sensorWorker.Run(nil)
	if errWorker != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	// Execute the workflow as if this were another client
	params := activity.WizSensorParams{
		KubeconfigPath:          "/path/to/kubeconfig",
		ImagePullSecretUsername: "your-username",
		ImagePullSecretPassword: "your-password",
		WizApiTokenClientId:     "your-client-id",
		WizApiTokenClientToken:  "your-client-token",
	}

	we, errEx := temporalCl.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        "deploy-wiz-sensor-workflow-001",
		TaskQueue: "wiz-sensor-task-queue",
	}, workflow.DeploySensorWorkflow, params)
	if errEx != nil {
		log.Fatalf("Failed to execute workflow: %v", errEx)
	}

	log.Printf("Workflow started successfully. Workflow ID: %s, Run ID: %s", we.GetID(), we.GetRunID())
}
