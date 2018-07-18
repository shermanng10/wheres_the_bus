package main

import (
	"testing"
)

func TestNewAlexaTextResponse(t *testing.T) {
	t.Run("Returns proper structure", func(t *testing.T) {
		resp := NewAlexaTextResponse("Selim Bradley is A Homonculus.")

		expectedVersion := "1.0"
		expectedOutputSpeechType := "PlainText"
		expectedOutputSpeechText := "Selim Bradley is A Homonculus."

		if expectedVersion != resp.Version {
			t.Errorf("Expected version %v, got %v", expectedVersion, resp.Version)
		}

		if expectedOutputSpeechType != resp.Response.OutputSpeech.Type {
			t.Errorf("Expected output speech type %v, got %v", expectedOutputSpeechType, resp.Response.OutputSpeech.Type)
		}

		if expectedOutputSpeechText != resp.Response.OutputSpeech.Text {
			t.Errorf("Expected output speech text %v, got %v", expectedOutputSpeechText, resp.Response.OutputSpeech.Text)
		}
	})
}
