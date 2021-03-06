package workflows

import (
	"time"

	"go.temporal.io/temporal"
	"go.temporal.io/temporal/workflow"

	"github.com/temporalio/temporal-go-demo/common"
)

type (
	AccountTransferRequest struct {
		FromAccountId string
		ToAccountId   string
		ReferenceId   string
		Amount        int
	}
)

func TransferWorkflow(ctx workflow.Context, transferRequest AccountTransferRequest) error {
	ao := workflow.ActivityOptions{
		TaskList:               common.ActivityTaskList,
		ScheduleToStartTimeout: 10 * time.Minute,
		StartToCloseTimeout:    5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			MaximumInterval:    10 * time.Second,
			ExpirationInterval: 10 * time.Minute,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	logger.Info("transfer workflow started")

	err := workflow.ExecuteActivity(ctx, "withdraw",
		transferRequest.FromAccountId,
		transferRequest.ReferenceId,
		transferRequest.Amount).Get(ctx, nil)
	if err != nil {
		return err
	}
	logger.Info("withdrawal completed")

	err = workflow.ExecuteActivity(ctx, "deposit",
		transferRequest.ToAccountId,
		transferRequest.ReferenceId,
		transferRequest.Amount).Get(ctx, nil)
	if err != nil {
		return err
	}
	logger.Info("deposit completed")

	return nil
}
