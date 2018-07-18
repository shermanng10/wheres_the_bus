package main

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestBusLocationHandler(t *testing.T) {
	t.Run("Test Gets Expected Response (No Bus Times)", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockBusService := NewMockBusLocationService(mockCtrl)
		mockBusService.EXPECT().GetBusTimesByStopCode("503471").Return([]BusTime{}, nil).Times(1)

		busHandler := NewBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes("503471")
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There are no buses arriving at the stop."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}

		expectedOutputSpeechType := "PlainText"
		actualOutputSpeechType := resp.Response.OutputSpeech.Type
		if actualOutputSpeechType != expectedOutputSpeechType {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechType, actualOutputSpeechType)
		}
	})

	t.Run("Test Gets Expected Response (1 Bus Time)", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
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

		busHandler := NewBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes("503471")
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There is one bus coming, the Q59 which is 1 stop away."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}

		expectedOutputSpeechType := "PlainText"
		actualOutputSpeechType := resp.Response.OutputSpeech.Type
		if actualOutputSpeechType != expectedOutputSpeechType {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechType, actualOutputSpeechType)
		}
	})

	t.Run("Test Gets Expected Response (More Than One Bus Time)", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
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

		busHandler := NewBusLocationHandler(mockBusService)
		resp, err := busHandler.GetBusTimes("503471")
		if err != nil {
			t.Errorf("Did not expect an error, got %v", err)
		}

		expectedOutputSpeechText := "There are 3 buses heading toward the stop, the Q59 which is 1 stop away, the Q58 which is 1.3 miles away, the Q58 which is 1.9 miles away."
		actualOutputSpeechText := resp.Response.OutputSpeech.Text
		if actualOutputSpeechText != expectedOutputSpeechText {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechText, actualOutputSpeechText)
		}

		expectedOutputSpeechType := "PlainText"
		actualOutputSpeechType := resp.Response.OutputSpeech.Type
		if actualOutputSpeechType != expectedOutputSpeechType {
			t.Errorf("Expected %v, got: %v", expectedOutputSpeechType, actualOutputSpeechType)
		}
	})
}
