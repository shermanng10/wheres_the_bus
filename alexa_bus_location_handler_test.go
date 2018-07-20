package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestBusLocationHandlerGetBusTimes(t *testing.T) {
	testTable := []struct {
		testName                 string
		mockBusStopCode          string
		mockResponse             []BusTime
		mockError                error
		expectedOutputSpeechText string
		alexaRequest             AlexaRequest
		alexaContext             context.Context
		err                      error
	}{

		{
			testName:                 "Gets Expected Response (No Bus Times)",
			mockBusStopCode:          "503471",
			mockResponse:             []BusTime{},
			mockError:                nil,
			expectedOutputSpeechText: "There are no buses arriving at the stop.",
			alexaRequest:             AlexaRequest{},
			alexaContext:             nil,
			err:                      nil,
		},

		{
			testName:        "Gets Expected Response (1 Bus Time) When Arrival Time is 0 Mins Away",
			mockBusStopCode: "503471",
			mockResponse: []BusTime{
				{
					Stop:          "503471",
					BusName:       "Q59",
					Distance:      "1 stop away",
					ArrivalTime:   time.Now(),
					DepartureTime: time.Now(),
				},
			},
			mockError:                nil,
			expectedOutputSpeechText: "There is one bus coming, the Q59 which is 1 stop away.",
			alexaRequest:             AlexaRequest{},
			alexaContext:             nil,
			err:                      nil,
		},

		{
			testName:        "Gets Expected Response (More Than One Bus Time) When Arrival Time is 0 Mins Away",
			mockBusStopCode: "503471",
			mockResponse: []BusTime{
				{
					Stop:          "503471",
					BusName:       "Q59",
					Distance:      "1 stop away",
					ArrivalTime:   time.Now(),
					DepartureTime: time.Now(),
				},
				{
					Stop:          "503471",
					BusName:       "Q58",
					Distance:      "1.3 miles away",
					ArrivalTime:   time.Now(),
					DepartureTime: time.Now(),
				},
				{
					Stop:          "503471",
					BusName:       "Q58",
					Distance:      "1.9 miles away",
					ArrivalTime:   time.Now(),
					DepartureTime: time.Now(),
				},
			},
			mockError:                nil,
			expectedOutputSpeechText: "There are 3 buses heading toward the stop, the Q59 which is 1 stop away, the Q58 which is 1.3 miles away, the Q58 which is 1.9 miles away.",
			alexaRequest:             AlexaRequest{},
			alexaContext:             nil,
			err:                      nil,
		},

		{
			testName:        "Gets Expected Response (1 Bus Time) When Arrival Time is > 0 Mins Away",
			mockBusStopCode: "503471",
			mockResponse: []BusTime{
				{
					Stop:          "503471",
					BusName:       "Q59",
					Distance:      "1 stop away",
					ArrivalTime:   time.Now().Local().Add(time.Duration(120) * time.Minute),
					DepartureTime: time.Now().Local().Add(time.Duration(120) * time.Minute),
					MinsAway:      120,
				},
			},
			mockError:                nil,
			expectedOutputSpeechText: "There is one bus coming, the Q59 which is 120 minutes away.",
			alexaRequest:             AlexaRequest{},
			alexaContext:             nil,
			err:                      nil,
		},

		{
			testName:        "Gets Expected Response (More Than One Bus Time) When Arrival Time is > 0 Mins Away",
			mockBusStopCode: "503471",
			mockResponse: []BusTime{
				{
					Stop:          "503471",
					BusName:       "Q59",
					Distance:      "1 stop away",
					ArrivalTime:   time.Now().Local().Add(time.Duration(7) * time.Minute),
					DepartureTime: time.Now().Local().Add(time.Duration(7) * time.Minute),
					MinsAway:      7,
				},
				{
					Stop:          "503471",
					BusName:       "Q58",
					Distance:      "1.3 miles away",
					ArrivalTime:   time.Now().Local().Add(time.Duration(10) * time.Minute),
					DepartureTime: time.Now().Local().Add(time.Duration(10) * time.Minute),
					MinsAway:      10,
				},
				{
					Stop:          "503471",
					BusName:       "Q58",
					Distance:      "1.9 miles away",
					ArrivalTime:   time.Now().Local().Add(time.Duration(43) * time.Minute),
					DepartureTime: time.Now().Local().Add(time.Duration(43) * time.Minute),
					MinsAway:      43,
				},
			},
			mockError:                nil,
			expectedOutputSpeechText: "There are 3 buses heading toward the stop, the Q59 which is 7 minutes away, the Q58 which is 10 minutes away, the Q58 which is 43 minutes away.",
			alexaRequest:             AlexaRequest{},
			alexaContext:             nil,
			err:                      nil,
		},

		{
			testName:                 "Gets Expected Error",
			mockBusStopCode:          "503471",
			mockResponse:             []BusTime{},
			mockError:                errors.New("Unexpected error"),
			expectedOutputSpeechText: "",
			alexaRequest:             AlexaRequest{},
			alexaContext:             nil,
			err:                      errors.New("Unexpected error"),
		},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockBusService := NewMockBusLocationService(mockCtrl)
			mockBusService.EXPECT().GetBusTimesByStopCode(test.mockBusStopCode).Return(
				test.mockResponse,
				test.mockError,
			).Times(1)

			busHandler := NewAlexaBusLocationHandler(mockBusService)
			resp, err := busHandler.GetBusTimes(test.alexaContext, test.alexaRequest)
			if err != nil && test.err != nil && err.Error() != test.err.Error() {
				t.Errorf("Expected %v error, got %v", test.err, err)
			}

			actualOutputSpeechText := resp.Response.OutputSpeech.Text
			if actualOutputSpeechText != test.expectedOutputSpeechText {
				t.Errorf("Expected %v, got: %v", test.expectedOutputSpeechText, actualOutputSpeechText)
			}
		})

	}
}
