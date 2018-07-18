package main

import (
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
	t.Run("Test Handles Expected Valid Data", func(t *testing.T) {
		stubResponse := fixture("standard_response.json")

		var stubHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(stubResponse))
		})

		var stubServer = httptest.NewServer(stubHandler)
		defer stubServer.Close()

		api := NewMTAStopMonitoringAPI(&http.Client{})
		api.SetBaseUrl(stubServer.URL)

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

	t.Run("Test Is Not Catastrophic On Broken Data (Bad JSON Keys or Structure)", func(t *testing.T) {
		stubResponse := fixture("bad_response_with_wonky_keys.json")

		var stubHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(stubResponse))
		})

		var stubServer = httptest.NewServer(stubHandler)
		defer stubServer.Close()

		api := NewMTAStopMonitoringAPI(&http.Client{})
		api.SetBaseUrl(stubServer.URL)

		busTimes, err := api.GetBusTimesByStopCode("503471")
		if err != nil {
			t.Logf("Did not expect an error %v", err)
			t.FailNow()
		}

		expectedNumberOfBusTimes := 0
		if len(busTimes) != expectedNumberOfBusTimes {
			t.Errorf("Expected %v bus times got: %v", expectedNumberOfBusTimes, len(busTimes))
		}
	})
}
