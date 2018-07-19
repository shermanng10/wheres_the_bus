package main

import (
	"context"
	"fmt"
	"strings"
)

type BusLocationHandler interface {
	GetBusTimes(context.Context, AlexaRequest) (AlexaTextResponse, error)
}

func AlexaBusLocationHandlerFactory() *AlexaBusLocationHandler {
	return NewAlexaBusLocationHandler(MTAStopMonitoringAPIFactory())
}

func NewAlexaBusLocationHandler(busService BusLocationService) *AlexaBusLocationHandler {
	return &AlexaBusLocationHandler{
		busService: busService,
	}
}

type AlexaBusLocationHandler struct {
	busService BusLocationService
}

func (h *AlexaBusLocationHandler) GetBusTimes(ctx context.Context, r AlexaRequest) (AlexaTextResponse, error) {
	var stopCode string
	slot, _ := r.Request.Intent.Slots["stopCode"]
	if slot.Value != "" {
		stopCode = slot.Value
	} else {
		stopCode = "503471" // Hard coded to my stop for now will be dynamic to saved preference in future.
	}

	busTimes, err := h.busService.GetBusTimesByStopCode(stopCode)
	if err != nil {
		return AlexaTextResponse{}, err
	}

	return h.makeBusTimesResponse(busTimes), nil
}

func (h *AlexaBusLocationHandler) makeBusTimesResponse(busTimes []BusTime) AlexaTextResponse {
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
