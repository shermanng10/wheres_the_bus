package main

import (
	"context"
	"fmt"
	"strings"
)

type BusLocationHandler interface {
	GetBusTimes(context.Context, AlexaRequest) (AlexaTextResponse, error)
}

func InitAlexaBusLocationHandler() *AlexaBusLocationHandler {
	return NewAlexaBusLocationHandler(InitMTAStopMonitoringAPI(), InitDynamoBusStopPreferenceStore())
}

func NewAlexaBusLocationHandler(busService BusLocationService, prefStore BusStopPreferenceStore) *AlexaBusLocationHandler {
	return &AlexaBusLocationHandler{
		busService:      busService,
		preferenceStore: prefStore,
	}
}

type AlexaBusLocationHandler struct {
	busService      BusLocationService
	preferenceStore BusStopPreferenceStore
}

func (h *AlexaBusLocationHandler) GetBusTimes(ctx context.Context, r AlexaRequest) (AlexaTextResponse, error) {
	var stopCode string
	stopCodeSlot, _ := r.Request.Intent.Slots["stopCode"]

	if stopCodeSlot.Value != "" {
		stopCode = stopCodeSlot.Value
	} else {
		var err error
		stopCode, err = h.preferenceStore.GetStopCodePreference(r.Session.User.UserId, "default")
		if err != nil {
			return AlexaTextResponse{}, err
		}
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
		var message string
		// If the bus is zero minutes away we want to use the mile distance so it sounds better, other wise we
		// use the minute distance because it is more accurate.
		if busTime.MinsAway == 0 {
			message = fmt.Sprintf(" the %v which is %v", busTime.BusName, busTime.Distance)
		} else {
			message = fmt.Sprintf(" the %v which is %v minutes away", busTime.BusName, busTime.MinsAway)
		}

		busMessageStrings = append(busMessageStrings, message)
	}

	outputMessage += fmt.Sprintf("%v.", strings.Join(busMessageStrings, ","))

	return NewAlexaTextResponse(outputMessage)
}
