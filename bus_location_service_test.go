package main

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestGetBusTimesByStopCode(t *testing.T) {
	t.Run("Test Expected Valid Data", func(t *testing.T) {
		mockResponse := fixture("standard_response.json")

		var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		})

		var mockServer = httptest.NewServer(mockHandler)
		defer mockServer.Close()

		api := NewMTAStopMonitoringAPI(&http.Client{})
		api.SetBaseUrl(mockServer.URL)

		busTimes, err := api.GetBusTimesByStopCode("503471")
		if err != nil {
			t.Logf("Did not expect an error %v", err)
			t.FailNow()
		}

		expectedNumberOfBusTimes := 3
		if len(busTimes) != expectedNumberOfBusTimes {
			t.Errorf("Expected %v bus times got: %v", expectedNumberOfBusTimes, len(busTimes))
		}

		expectedStopCode := "503471"
		for _, busTime := range busTimes {
			if busTime.Stop != expectedStopCode {
				t.Errorf("Expected %v bus code, got: %v", expectedStopCode, busTime.Stop)
			}

			err := busTime.Validate()
			if err != nil {
				t.Errorf("Did not expect an error, got: %v", err)
			}
		}
	})
}
