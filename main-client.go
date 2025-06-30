package main

import (
	"context"
	"github.com/jefferyfry/temporal-wiz-sensor/activity"
	"github.com/jefferyfry/temporal-wiz-sensor/workflow"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	temporalCl, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer temporalCl.Close()

	// Execute the workflow as if this were another client
	params := activity.WizSensorParams{
		KubeconfigPath:          "<>", // "/Users/jdoe/.kube/config",
		KubeconfigContext:       "<>", //"arn:aws:eks:us-east-1:12346789:cluster/agent-app-cluster",
		ImagePullSecretUsername: "<>", // "wizio-repo-.....",
		ImagePullSecretPassword: "<>", // "CY6eLf.....",
		WizApiTokenClientId:     "<>", // "t3zxgy4libc2jfi.....",
		WizApiTokenClientToken:  "<>", // "iLEuEe7pw9T5u3GNivxKt1Oc1vYnI.....",
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
