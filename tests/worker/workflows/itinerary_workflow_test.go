package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"go.uber.org/cadence/testsuite"

	"training/hello-cadence/app/itinerary_cadence/workflows"
	"training/hello-cadence/app/models"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_SimpleWorkflow() {
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(workflows.ArrivalPointSignal, models.Point{
			Lat: 48.767999, Long: 2.298879,
		})
	}, 0)

	s.env.ExecuteWorkflow(workflows.ComputeItineraryWorkflow, models.Point{
		Lat: 46.198941, Long: 6.140618,
	})

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var res models.Itinerary
	s.env.GetWorkflowResult(&res)
	s.Equal(res.Road[0].Lat, float32(47.483467))
	s.Equal(res.Road[0].Long, float32(4.2197485))
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
