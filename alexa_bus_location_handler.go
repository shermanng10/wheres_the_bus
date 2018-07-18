package main

import (
	"fmt"

	"strings"
)

type AlexaBusLocationHandler interface {
	GetBusTimes(string) (AlexaTextResponse, error)
}

func NewBusLocationHandler(busService BusLocationService) *BusLocationHandler {
	return &BusLocationHandler{
		busService: busService,
	}
}

type BusLocationHandler struct {
	busService BusLocationService
}

func (h *BusLocationHandler) GetBusTimes(stopCode string) (AlexaTextResponse, error) {
	busTimes, err := h.busService.GetBusTimesByStopCode(stopCode)
	if err != nil {
		return NewAlexaTextResponse("Something went wrong with your request."), err
	}

	return h.makeBusTimesResponse(busTimes), nil
}

func (h *BusLocationHandler) makeBusTimesResponse(busTimes []BusTime) AlexaTextResponse {
	var outputMessage string
	numTimes := len(busTimes)

	if numTimes == 0 {
		outputMessage = "There are no buses arriving at the stop"
	} else if numTimes == 1 {
		outputMessage = "There is one bus coming,"
	} else {
		outputMessage = fmt.Sprintf("There are %v buses heading toward the stop,", numTimes)
	}

	var busMessageStrings []string
	for _, busTime := range busTimes {
		busMessageStrings = append(busMessageStrings, fmt.Sprintf(" the %v which is %v", busTime.BusName, busTime.Distance))
	}

	outputMessage += fmt.Sprintf("%v.", strings.Join(busMessageStrings, ","))

	return NewAlexaTextResponse(outputMessage)
}
