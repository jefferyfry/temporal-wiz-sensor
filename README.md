# Temporal Wiz Sensor

Reusable Temproral Wiz Sensor install activity for k8s to be used in your Temporal workflows.

## Features
- Install and validate Temporal Wiz Sensor on Kubernetes clusters.
- Can be used in Temporal workflows to ensure the sensor is installed before proceeding with other activities.
- Packaged as a Go Module

## Installation
To use the Temporal Wiz Sensor in your Temporal workflows, you can install it as a Go module. Run the following command in your project directory:

```bash
go get github.com/jefferyfry/temporal-wiz-sensor@v0.1.0
```

## Usage
To use the Temporal Wiz Sensor in your Temporal workflows, you can import the package and use the `InstallSensor` activity. Here's an example of how to do this:

```go   
package workflows

import (
	"github.com/jefferyfry/temporal-wiz-sensor/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func DeploySensorWorkflow(ctx workflow.Context, params activity.WizSensorParams) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	return workflow.ExecuteActivity(ctx, activity.InstallWizSensorActivity, params).Get(ctx, nil)
}
```