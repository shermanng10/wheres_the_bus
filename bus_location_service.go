package main

import (
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"

	// Using this library for parsing deeply nested json so that I do not have to define many structs within structs
	"github.com/tidwall/gjson"
)

type BusLocationService interface {
	GetBusTimesByStopCode(string) ([]BusTime, error)
}

func MTAStopMonitoringAPIFactory() *MTABusStopMonitoringAPI {
	const defaultTimeOut = time.Second * 15
	defaultHttpClient := &http.Client{
		Timeout: defaultTimeOut,
	}
	return NewMTAStopMonitoringAPI(defaultHttpClient)
}

func NewMTAStopMonitoringAPI(httpClient *http.Client) *MTABusStopMonitoringAPI {
	const mtaStopMonitoringApiUrl = "http://bustime.mta.info/api/siri/stop-monitoring.json"
	return &MTABusStopMonitoringAPI{
		HttpClient: httpClient,
		baseUrl:    mtaStopMonitoringApiUrl,
	}
}

type MTABusStopMonitoringAPI struct {
	baseUrl    string
	HttpClient *http.Client
}

type BusTime struct {
	Stop          string
	ArrivalTime   time.Time
	DepartureTime time.Time
	MinsAway      float64
	Distance      string
	BusName       string
}

func (api *MTABusStopMonitoringAPI) SetBaseUrl(url string) {
	api.baseUrl = url
}

func (api *MTABusStopMonitoringAPI) GetBusTimesByStopCode(code string) ([]BusTime, error) {
	endpoint := api.makeStopMonitoringEndpoint(code)
	resp, err := api.HttpClient.Get(endpoint)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return api.responseToBusTimes(body, code), nil
}

func (api *MTABusStopMonitoringAPI) makeStopMonitoringEndpoint(stopCode string) string {
	const detailLevel = "minimum"
	const operatorRef = "MTA"
	const apiVersion = "2"
	const maxStopResults = "3"

	form := url.Values{}
	form.Set("key", os.Getenv("MTA_STOP_MONITORING_API_KEY"))
	form.Set("version", apiVersion)
	form.Set("OperatorRef", operatorRef)
	form.Set("MonitoringRef", stopCode)
	form.Set("StopMonitoringDetailLevel", detailLevel)
	form.Set("MaximumStopVisits", maxStopResults)
	qs := form.Encode()
	return api.baseUrl + "?" + qs
}

func (api *MTABusStopMonitoringAPI) responseToBusTimes(resp []byte, stopCode string) []BusTime {
	const stopsVisitJsonPath = "Siri.ServiceDelivery.StopMonitoringDelivery.0.MonitoredStopVisit"
	// All constant paths that follow are nested within stopsVisitJsonPath so that we don't have to re-parse
	// the entire object multiple times.
	const busNameJsonPath = "MonitoredVehicleJourney.PublishedLineName"
	const arrivalTimeJsonPath = "MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime"
	const departureTimeJsonPath = "MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime"
	const distanceJsonPath = "MonitoredVehicleJourney.MonitoredCall.ArrivalProximityText"

	var busTimes []BusTime
	stops := gjson.GetBytes(resp, stopsVisitJsonPath)
	for _, stop := range stops.Array() {
		stopJson := stop.String()
		arrivalTime := gjson.Get(stopJson, arrivalTimeJsonPath).Time()

		busTimes = append(busTimes, BusTime{
			Stop:          stopCode,
			ArrivalTime:   arrivalTime,
			DepartureTime: gjson.Get(stopJson, departureTimeJsonPath).Time(),
			Distance:      gjson.Get(stopJson, distanceJsonPath).String(),
			BusName:       gjson.Get(stopJson, busNameJsonPath).String(),
			MinsAway:      api.calculateMinsAway(arrivalTime),
		})
	}

	return busTimes
}

func (api *MTABusStopMonitoringAPI) calculateMinsAway(arrivalTime time.Time) float64 {
	now := time.Now()
	if arrivalTime.IsZero() || now.After(arrivalTime) {
		// Return a zero valued time away if the arrival time is unknown or the current time is after the arrival so
		// that in the response, we can choose if we want to include the time based on if it is accurate.
		return 0
	} else {
		duration := time.Until(arrivalTime)
		return math.Ceil(duration.Minutes())
	}
}
