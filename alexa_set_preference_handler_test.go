package main

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestAlexaSetPreferenceHandler(t *testing.T) {
	testTable := []struct {
		mockBusStopPreferenceStoreUserId      string
		mockBusStopPreferenceStorePrefName    string
		mockBusStopPreferenceStoreStopCode    string
		mockBusStopPreferenceStoreCalledTimes int
		mockBusStopPreferenceStoreResponse    string
		mockBusStopPreferenceStoreError       error
		testName                              string
		alexaRequest                          AlexaRequest
		expectedOutputSpeechText              string
		expectedError                         error
	}{

		{
			testName: "Returns proper error response when stop code not provided.",
			mockBusStopPreferenceStoreUserId:      "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			mockBusStopPreferenceStorePrefName:    "default",
			mockBusStopPreferenceStoreStopCode:    "",
			mockBusStopPreferenceStoreError:       nil,
			mockBusStopPreferenceStoreCalledTimes: 0,
			alexaRequest:                          newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", ""),
			expectedOutputSpeechText:              "",
			expectedError:                         errors.New("Sorry, I didn't hear a stop code."),
		},

		{
			testName: "Returns proper error response when user id is missing.",
			mockBusStopPreferenceStoreUserId:      "",
			mockBusStopPreferenceStorePrefName:    "default",
			mockBusStopPreferenceStoreStopCode:    "503471",
			mockBusStopPreferenceStoreError:       nil,
			mockBusStopPreferenceStoreCalledTimes: 0,
			alexaRequest:                          newFakeAlexaRequest("", "503471"),
			expectedOutputSpeechText:              "",
			expectedError:                         errors.New("Something went wrong, user id not given."),
		},

		{
			testName: "Returns proper response on successful save.",
			mockBusStopPreferenceStoreUserId:      "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			mockBusStopPreferenceStorePrefName:    "default",
			mockBusStopPreferenceStoreStopCode:    "503471",
			mockBusStopPreferenceStoreError:       nil,
			mockBusStopPreferenceStoreCalledTimes: 1,
			alexaRequest:                          newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "503471"),
			expectedOutputSpeechText:              "Saved your default bus stop code to 503471.",
			expectedError:                         nil,
		},

		{
			testName: "Returns proper error response on unexpected error.",
			mockBusStopPreferenceStoreUserId:      "amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT",
			mockBusStopPreferenceStorePrefName:    "default",
			mockBusStopPreferenceStoreStopCode:    "503471",
			mockBusStopPreferenceStoreError:       errors.New("Unexpected error."),
			mockBusStopPreferenceStoreCalledTimes: 1,
			alexaRequest:                          newFakeAlexaRequest("amzn1.ask.account.amzn1.ask.account.FAKEACCOUNT", "503471"),
			expectedOutputSpeechText:              "",
			expectedError:                         errors.New("Unexpected error."),
		},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockBusDataStore := NewMockBusStopPreferenceStore(mockCtrl)
			mockBusDataStore.EXPECT().SetStopCodePreference(
				test.mockBusStopPreferenceStoreUserId,
				test.mockBusStopPreferenceStorePrefName,
				test.mockBusStopPreferenceStoreStopCode,
			).Return(
				test.mockBusStopPreferenceStoreError,
			).Times(test.mockBusStopPreferenceStoreCalledTimes)

			prefHandler := NewAlexaSetPreferenceHandler(mockBusDataStore)
			resp, err := prefHandler.SetBusPreference(nil, test.alexaRequest)
			if err != nil && test.expectedError != nil && err.Error() != test.expectedError.Error() {
				t.Errorf("Expected %v error, got %v", test.expectedError, err)
			}

			actualOutputSpeechText := resp.Response.OutputSpeech.Text
			if actualOutputSpeechText != test.expectedOutputSpeechText {
				t.Errorf("Expected %v, got: %v", test.expectedOutputSpeechText, actualOutputSpeechText)
			}
		})
	}

}
