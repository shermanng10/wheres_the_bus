package main

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestBusLocationHandler(t *testing.T) {
	t.Run("Test Gets Expected Response (No Bus Times)", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBusService := NewMockBusLocationService(mockCtrl)
		mockBusService.EXPECT().GetBusTimesByStopCode("503471").Return([]BusTime{}, nil).Times(1)

		busHandler := NewAlexaBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes(nil, AlexaRequest{})
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There are no buses arriving at the stop."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}
	})

	t.Run("Test Gets Expected Response (1 Bus Time) When Arrival Time is 0 Mins Away", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBusService := NewMockBusLocationService(mockCtrl)
		mockBusService.EXPECT().GetBusTimesByStopCode("503471").Return([]BusTime{
			{
				Stop:          "503471",
				BusName:       "Q59",
				Distance:      "1 stop away",
				ArrivalTime:   time.Now(),
				DepartureTime: time.Now(),
			},
		}, nil).Times(1)

		busHandler := NewAlexaBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes(nil, AlexaRequest{})
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There is one bus coming, the Q59 which is 1 stop away."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}
	})

	t.Run("Test Gets Expected Response (More Than One Bus Time) When Arrival Time is 0 Mins Away", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBusService := NewMockBusLocationService(mockCtrl)
		mockBusService.EXPECT().GetBusTimesByStopCode("503471").Return([]BusTime{
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
		}, nil).Times(1)

		busHandler := NewAlexaBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes(nil, AlexaRequest{})
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There are 3 buses heading toward the stop, the Q59 which is 1 stop away, the Q58 which is 1.3 miles away, the Q58 which is 1.9 miles away."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}
	})

	t.Run("Test Gets Expected Response (1 Bus Time) When Arrival Time is > 0 Mins Away", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBusService := NewMockBusLocationService(mockCtrl)
		mockBusService.EXPECT().GetBusTimesByStopCode("503471").Return([]BusTime{
			{
				Stop:          "503471",
				BusName:       "Q59",
				Distance:      "1 stop away",
				ArrivalTime:   time.Now().Local().Add(time.Duration(120) * time.Minute),
				DepartureTime: time.Now().Local().Add(time.Duration(120) * time.Minute),
				MinsAway:      120,
			},
		}, nil).Times(1)

		busHandler := NewAlexaBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes(nil, AlexaRequest{})
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There is one bus coming, the Q59 which is 120 minutes away."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}
	})

	t.Run("Test Gets Expected Response (More Than One Bus Time) When Arrival Time is > 0 Mins Away", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBusService := NewMockBusLocationService(mockCtrl)
		mockBusService.EXPECT().GetBusTimesByStopCode("503471").Return([]BusTime{
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
		}, nil).Times(1)

		busHandler := NewAlexaBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes(nil, AlexaRequest{})
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There are 3 buses heading toward the stop, the Q59 which is 7 minutes away, the Q58 which is 10 minutes away, the Q58 which is 43 minutes away."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}
	})
}
