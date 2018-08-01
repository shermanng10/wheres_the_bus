package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func newFakeAlexaRequest(userId string, stopCode string, stopPrefName string) AlexaRequest {
	ar := AlexaRequest{}
	ar.Session.User.UserId = userId

	ar.Request.Intent.Slots = map[string]struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{}

	ar.Request.Intent.Slots["stopCode"] = struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  "fakeName",
		Value: stopCode,
	}

	ar.Request.Intent.Slots["stopPrefName"] = struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  "fakeName",
		Value: stopPrefName,
	}

	return ar
}

func TestBusLocationHandlerGetBusTimes(t *testing.T) {
	testTable := []struct {
		testName                               string
		expectedMockBusServiceCode             string
		mockBusServiceResponse                 []BusTime
		mockBusServiceError                    error
		expectedBusStopPreferenceStoreUserId   string
		expectedBusStopPreferenceStorePrefName string
		mockBusStopPreferenceStoreResponse     string
		mockBusStopPreferenceStoreError        error
		expectedOutputSpeechText               string
		alexaRequest                           AlexaRequest
		alexaContext                           context.Context
		err                                    error
	}{

		{
			testName:                               "Gets Expected Response (No Bus Times) when no stop pref name passed in request",
			expectedMockBusServiceCode:             "503471",
			mockBusServiceResponse:                 []BusTime{},
			mockBusServiceError:                    nil,
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "default",
			mockBusStopPreferenceStoreResponse:     "503471",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "There are no buses arriving at the stop.",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "", ""),
			alexaContext:                           nil,
			err:                                    nil,
		},

		{
			testName:                               "Gets Expected Response (No Bus Times) when stop pref name passed in request",
			expectedMockBusServiceCode:             "503471",
			mockBusServiceResponse:                 []BusTime{},
			mockBusServiceError:                    nil,
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "queens",
			mockBusStopPreferenceStoreResponse:     "503471",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "There are no buses arriving at the stop.",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "", "queens"),
			alexaContext:                           nil,
			err:                                    nil,
		},

		{
			testName:                   "Gets Expected Response (1 Bus Time) When Arrival Time is 0 Mins Away",
			expectedMockBusServiceCode: "503472",
			mockBusServiceResponse: []BusTime{
				{
					Stop:          "503472",
					BusName:       "Q59",
					Distance:      "1 stop away",
					ArrivalTime:   time.Now(),
					DepartureTime: time.Now(),
				},
			},
			mockBusServiceError:                    nil,
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "brooklyn",
			mockBusStopPreferenceStoreResponse:     "503472",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "There is one bus coming, the Q59 which is 1 stop away.",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "", "brooklyn"),
			alexaContext:                           nil,
			err:                                    nil,
		},

		{
			testName:                   "Gets Expected Response (More Than One Bus Time) When Arrival Time is 0 Mins Away",
			expectedMockBusServiceCode: "503471",
			mockBusServiceResponse: []BusTime{
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
			mockBusServiceError:                    nil,
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "manhattan",
			mockBusStopPreferenceStoreResponse:     "503471",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "There are 3 buses heading toward the stop, the Q59 which is 1 stop away, the Q58 which is 1.3 miles away, the Q58 which is 1.9 miles away.",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "", "manhattan"),
			alexaContext:                           nil,
			err:                                    nil,
		},

		{
			testName:                   "Gets Expected Response (1 Bus Time) When Arrival Time is > 0 Mins Away",
			expectedMockBusServiceCode: "503471",
			mockBusServiceResponse: []BusTime{
				{
					Stop:          "503471",
					BusName:       "Q59",
					Distance:      "1 stop away",
					ArrivalTime:   time.Now().Local().Add(time.Duration(120) * time.Minute),
					DepartureTime: time.Now().Local().Add(time.Duration(120) * time.Minute),
					MinsAway:      120,
				},
			},
			mockBusServiceError:                    nil,
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "default",
			mockBusStopPreferenceStoreResponse:     "503471",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "There is one bus coming, the Q59 which is 120 minutes away.",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "503471", ""),
			alexaContext:                           nil,
			err:                                    nil,
		},

		{
			testName:                   "Gets Expected Response (More Than One Bus Time) When Arrival Time is > 0 Mins Away",
			expectedMockBusServiceCode: "503471",
			mockBusServiceResponse: []BusTime{
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
			mockBusServiceError:                    nil,
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "default",
			mockBusStopPreferenceStoreResponse:     "503471",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "There are 3 buses heading toward the stop, the Q59 which is 7 minutes away, the Q58 which is 10 minutes away, the Q58 which is 43 minutes away.",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "503471", ""),
			alexaContext:                           nil,
			err:                                    nil,
		},

		{
			testName:                               "Gets Expected Error",
			expectedMockBusServiceCode:             "503471",
			mockBusServiceResponse:                 []BusTime{},
			mockBusServiceError:                    errors.New("Unexpected error."),
			expectedBusStopPreferenceStoreUserId:   "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			expectedBusStopPreferenceStorePrefName: "default",
			mockBusStopPreferenceStoreResponse:     "503471",
			mockBusStopPreferenceStoreError:        nil,
			expectedOutputSpeechText:               "",
			alexaRequest:                           newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "503471", ""),
			alexaContext:                           nil,
			err:                                    errors.New("Unexpected error."),
		},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockBusService := NewMockBusLocationService(mockCtrl)
			mockBusService.EXPECT().GetBusTimesByStopCode(test.expectedMockBusServiceCode).Return(
				test.mockBusServiceResponse,
				test.mockBusServiceError,
			).MaxTimes(1)

			mockBusDataStore := NewMockBusStopPreferenceStore(mockCtrl)
			mockBusDataStore.EXPECT().GetStopCodePreference(
				test.expectedBusStopPreferenceStoreUserId,
				test.expectedBusStopPreferenceStorePrefName,
			).Return(
				test.mockBusStopPreferenceStoreResponse,
				test.mockBusStopPreferenceStoreError,
			).MaxTimes(1)

			busHandler := NewAlexaBusLocationHandler(mockBusService, mockBusDataStore)
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
