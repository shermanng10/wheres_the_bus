package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
)

func makeTestIntentRequest(intent string) AlexaRequest {
	request := AlexaRequest{}
	request.Version = "1.0"
	request.Request.Type = "IntentRequest"
	request.Request.Time = "2018-07-18T05:42:26Z"
	request.Request.Intent.ConfirmationStatus = "NONE"
	request.Request.Intent.Name = intent
	return request
}

func TestAlexaMuxHandle(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	testTable := []struct {
		testName                       string
		expectedAlexaResponse          AlexaTextResponse
		expectedNumBusHandlerCalls     int
		expectedNumSetPrefHandlerCalls int
		testRequest                    AlexaRequest
	}{

		{
			testName:                       "wheres_the_bus intent calls bus handler with proper response",
			expectedAlexaResponse:          NewAlexaTextResponse("There is one bus coming, the Q59 which is 1 stop away."),
			expectedNumBusHandlerCalls:     1,
			expectedNumSetPrefHandlerCalls: 0,
			testRequest:                    makeTestIntentRequest("wheres_the_bus"),
		},

		{
			testName:                       "save_bus_code intent calls save preference handler",
			expectedAlexaResponse:          NewAlexaTextResponse("Saved your default bus stop code to 503471."),
			expectedNumBusHandlerCalls:     0,
			expectedNumSetPrefHandlerCalls: 1,
			testRequest:                    makeTestIntentRequest("save_bus_code"),
		},

		{
			testName:                       "AMAZON.HelpIntent intent calls with proper response",
			expectedAlexaResponse:          NewAlexaTextResponse("Just ask where the bus is."),
			expectedNumBusHandlerCalls:     0,
			expectedNumSetPrefHandlerCalls: 0,
			testRequest:                    makeTestIntentRequest("AMAZON.HelpIntent"),
		},

		{
			testName:                       "default intent calls with proper response",
			expectedAlexaResponse:          NewAlexaTextResponse("I'm sorry, I didn't understand that."),
			expectedNumBusHandlerCalls:     0,
			expectedNumSetPrefHandlerCalls: 0,
			testRequest:                    makeTestIntentRequest("default"),
		},
	}
	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockBusHandler := NewMockBusLocationHandler(mockCtrl)
			mockBusHandler.EXPECT().GetBusTimes(nil, test.testRequest).
				Return(test.expectedAlexaResponse, nil).
				Times(test.expectedNumBusHandlerCalls)

			mockSetPrefHandler := NewMockSetPreferenceHandler(mockCtrl)
			mockSetPrefHandler.EXPECT().SetBusPreference(nil, test.testRequest).
				Return(test.expectedAlexaResponse, nil).
				Times(test.expectedNumSetPrefHandlerCalls)

			alexaMux := NewAlexaMux(mockBusHandler, mockSetPrefHandler)
			resp, err := alexaMux.Handle(nil, test.testRequest)
			if err != nil {
				t.Errorf("Unexpected error, got: %v", err)
			}
			if resp != test.expectedAlexaResponse {
				t.Errorf("Expected response %v, got: %v", test.expectedAlexaResponse, resp)
			}
		})
	}
}
