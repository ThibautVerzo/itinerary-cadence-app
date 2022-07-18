package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/cadence/client"
	"go.uber.org/zap"

	"training/hello-cadence/app/adapters/cadenceAdapter"
	"training/hello-cadence/app/itinerary_cadence/workflows"
	"training/hello-cadence/app/models"
)

type ItineraryController struct {
	CadenceAdapter *cadenceAdapter.CadenceAdapter
	Logger         *zap.Logger
}

type workflowItineraryResponse struct {
	WorflowId, RunId, StatusURL string
}

// POST /itineraries.set-departure
// Initiates workflow itinerary computing with the given departure point.
// Returns the cadence WorkflowID and the RunID.
func (controller *ItineraryController) SendItineraryDeparture(ctx *gin.Context) {
	var departure models.Point
	ctx.BindJSON(&departure)

	workflowOptions := client.StartWorkflowOptions{
		TaskList:                     workflows.TaskListName,
		ExecutionStartToCloseTimeout: time.Hour * 24,
	}
	workflowExec, err := controller.CadenceAdapter.CadenceClient.StartWorkflow(
		context.Background(), workflowOptions, workflows.WorkFlowName, departure)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "Error starting workflow!")
		return
	}

	controller.Logger.Info(
		"Started work flow!", zap.String("WorkflowId", workflowExec.ID),
		zap.String("RunId", workflowExec.RunID),
	)

	ctx.IndentedJSON(http.StatusAccepted, workflowItineraryResponse{
		WorflowId: workflowExec.ID,
		RunId:     workflowExec.RunID,
		StatusURL: fmt.Sprintf("%s/itineraries.status?workflowId=%s",
			ctx.Request.Host, workflowExec.ID,
		),
	})
}

// POST /itineraries.set-arrival
// Sends signals to the given WorkflowID with the arrival point.
func (controller *ItineraryController) SendItineraryArrival(ctx *gin.Context) {
	workflowId := ctx.Query("workflowId")
	var arrival models.Point
	ctx.BindJSON(&arrival)

	err := controller.CadenceAdapter.CadenceClient.SignalWorkflow(
		context.Background(), workflowId, "", workflows.ArrivalPointSignal, arrival,
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, "Error while signaling: "+err.Error())
		return
	}

	controller.Logger.Info("Signaled work flow with the following params!",
		zap.String("WorkflowId", workflowId), zap.String("Arrival", arrival.String()),
	)

	ctx.IndentedJSON(http.StatusAccepted, workflowItineraryResponse{
		WorflowId: workflowId,
		StatusURL: fmt.Sprintf("%s/itineraries.status?workflowId=%s",
			ctx.Request.Host, workflowId),
	})
	return
}

// GET /itineraries.status
// Returns the itinerary computation status with the given WorkflowID.
func (controller *ItineraryController) GetItineraryStatus(ctx *gin.Context) {
	workflowId := ctx.Query("workflowId")

	workflow, err := controller.CadenceAdapter.CadenceClient.DescribeWorkflowExecution(
		context.Background(), workflowId, "",
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest,
			"Error while getting workflow status: "+err.Error(),
		)
		return
	}

	workflowCloseStatus := workflow.WorkflowExecutionInfo.CloseStatus
	if workflowCloseStatus == nil {
		ctx.IndentedJSON(http.StatusOK, "The workflow is running")
		return
	}
	ctx.IndentedJSON(
		http.StatusFound, "The workflow is close: "+workflowCloseStatus.String(),
	)
}

// GET /itineraries
// Returns the computed itinerary with the given WorkflowID.
func (controller *ItineraryController) GetItinerary(ctx *gin.Context) {
	workflowId := ctx.Query("workflowId")

	workflow := controller.CadenceAdapter.CadenceClient.GetWorkflow(
		context.Background(), workflowId, "",
	)

	var itinerary models.Itinerary
	workflow.Get(context.Background(), &itinerary)

	ctx.IndentedJSON(http.StatusOK, itinerary)
}
