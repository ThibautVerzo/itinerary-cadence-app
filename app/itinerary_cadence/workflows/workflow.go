package workflows

import (
	"context"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"

	"training/hello-cadence/app/gin_service/managers"
	"training/hello-cadence/app/models"
)

const TaskListName = "itineraryGroup"
const ArrivalPointSignal = "arrivalPointSignal"
const WorkFlowName = "itineraryWorkFlow"

func init() {
	var wo = workflow.RegisterOptions{
		Name: WorkFlowName,
	}
	workflow.RegisterWithOptions(ComputeItineraryWorkflow, wo)
	activity.Register(SetDepartureActivity)
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func SetDepartureActivity(ctx context.Context, departure models.Point) (models.Itinerary, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SetDepartureActivity activity started")

	// Here some computing can be done while waiting for the arrival point
	// (like compute near buses, tramway, train stops, ...).

	return models.Itinerary{
		Departure: &departure,
		Arrival:   nil,
		Road:      nil,
	}, nil
}

func ComputeItineraryWorkflow(ctx workflow.Context, departure models.Point) (*models.Itinerary, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("ComputeItineraryWorkflow workflow started")
	var itinerary models.Itinerary
	err := workflow.ExecuteActivity(ctx, SetDepartureActivity, departure).Get(
		ctx, &itinerary)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return nil, err
	}

	// After setting the departure point, the workflow will wait for the arrival point!
	selector := workflow.NewSelector(ctx)
	var arrival models.Point

	for {
		signalChan := workflow.GetSignalChannel(ctx, ArrivalPointSignal)
		selector.AddReceive(signalChan, func(ch workflow.Channel, more bool) {
			ch.Receive(ctx, &arrival)
			workflow.GetLogger(ctx).Info("Received arrival point from signal!",
				zap.String("signal", ArrivalPointSignal),
				zap.String("value", arrival.String()))
		})
		workflow.GetLogger(ctx).Info(
			"Waiting for signal on channel.. " + ArrivalPointSignal)
		// Wait for signal
		selector.Select(ctx)

		// Compute itinerary here
		itinerary.Arrival = &arrival
		managers.ComputeItinerary(&itinerary)
		return &itinerary, nil
	}
}
