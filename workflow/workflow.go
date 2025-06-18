package workflow

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
