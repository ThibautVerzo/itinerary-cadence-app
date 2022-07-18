package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"training/hello-cadence/app/adapters/cadenceAdapter"
	"training/hello-cadence/app/config"
	"training/hello-cadence/app/gin_service/controllers"
	"training/hello-cadence/app/itinerary_cadence/worker"
	"training/hello-cadence/app/itinerary_cadence/workflows"
	"training/hello-cadence/app/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()

	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	itineraryController := controllers.ItineraryController{
		CadenceAdapter: &cadenceClient,
		Logger:         appConfig.Logger,
	}

	router.POST("/itineraries.set-departure", itineraryController.SendItineraryDeparture)
	router.POST("/itineraries.set-arrival", itineraryController.SendItineraryArrival)
	router.GET("/itineraries.status", itineraryController.GetItineraryStatus)
	router.GET("/itineraries", itineraryController.GetItinerary)

	worker.StartWorkers(&cadenceClient, workflows.TaskListName)

	return router
}

type workflowItineraryResponse struct {
	WorflowId, RunId, StatusURL string
}

func TestItineraryComputing_Success(t *testing.T) {
	router := SetUpRouter()

	// Set new departure point (Initiate itinerary computimg workflow)
	httpRecorder := httptest.NewRecorder()
	departure := models.Point{Lat: 46.198941, Long: 6.140618}
	departureJSON, _ := json.Marshal(departure)
	req, _ := http.NewRequest(
		"POST", "/itineraries.set-departure", bytes.NewBuffer(departureJSON),
	)
	router.ServeHTTP(httpRecorder, req)
	assert.Equal(t, http.StatusAccepted, httpRecorder.Code)
	var workflowInfo workflowItineraryResponse
	json.Unmarshal(httpRecorder.Body.Bytes(), &workflowInfo)

	// Check if workflow is running
	httpRecorderStatus1 := httptest.NewRecorder()
	reqStatus1, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/itineraries.status?workflowId=%s", workflowInfo.WorflowId),
		nil,
	)
	router.ServeHTTP(httpRecorderStatus1, reqStatus1)
	assert.Equal(t, http.StatusOK, httpRecorderStatus1.Code)
	assert.Equal(t, "\"The workflow is running\"", httpRecorderStatus1.Body.String())

	// Set arrival point
	httpRecorderArrival := httptest.NewRecorder()
	arrival := models.Point{Lat: 48.767999, Long: 2.298879}
	arrivalJSON, _ := json.Marshal(arrival)
	reqArrival, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("/itineraries.set-arrival?workflowId=%s", workflowInfo.WorflowId),
		bytes.NewBuffer(arrivalJSON),
	)
	router.ServeHTTP(httpRecorderArrival, reqArrival)
	assert.Equal(t, http.StatusAccepted, httpRecorderArrival.Code)

	time.Sleep(2 * time.Second)

	// Check if workflow is completed
	httpRecorderStatus2 := httptest.NewRecorder()
	reqStatus2, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/itineraries.status?workflowId=%s", workflowInfo.WorflowId),
		nil,
	)
	router.ServeHTTP(httpRecorderStatus2, reqStatus2)
	assert.Equal(t, http.StatusFound, httpRecorderStatus2.Code)
	assert.Equal(t, "\"The workflow is close: COMPLETED\"", httpRecorderStatus2.Body.String())

	// Get itinerary result
	httpRecorderRes := httptest.NewRecorder()
	reqRes, _ := http.NewRequest(
		"GET", fmt.Sprintf("/itineraries?workflowId=%s", workflowInfo.WorflowId), nil,
	)
	router.ServeHTTP(httpRecorderRes, reqRes)
	assert.Equal(t, http.StatusOK, httpRecorderRes.Code)
	var res models.Itinerary
	json.Unmarshal(httpRecorderRes.Body.Bytes(), &res)
	assert.Equal(t, res.Road[0].Lat, float32(47.483467))
	assert.Equal(t, res.Road[0].Long, float32(4.2197485))

}
