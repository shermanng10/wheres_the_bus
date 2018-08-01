package main

import (
	"context"
	"errors"
	"fmt"
)

type SetPreferenceHandler interface {
	SetBusPreference(context.Context, AlexaRequest) (AlexaTextResponse, error)
}

func InitAlexaSetPreferenceHandler() *AlexaSetPreferenceHandler {
	return NewAlexaSetPreferenceHandler(InitDynamoBusStopPreferenceStore())
}

func NewAlexaSetPreferenceHandler(prefStore BusStopPreferenceStore) *AlexaSetPreferenceHandler {
	return &AlexaSetPreferenceHandler{
		preferenceStore: prefStore,
	}
}

type AlexaSetPreferenceHandler struct {
	preferenceStore BusStopPreferenceStore
}

func (h *AlexaSetPreferenceHandler) SetBusPreference(c context.Context, r AlexaRequest) (AlexaTextResponse, error) {
	userId := r.Session.User.UserId
	stopCodeSlot, _ := r.Request.Intent.Slots["stopCode"]
	if stopCodeSlot.Value == "" {
		return AlexaTextResponse{}, errors.New("Sorry, I didn't hear a stop code.")
	}
	if userId == "" {
		return AlexaTextResponse{}, errors.New("Something went wrong, user id not given.")
	}

	stopCode := stopCodeSlot.Value
	err := h.preferenceStore.SetStopCodePreference(userId, "default", stopCode)
	if err != nil {
		return AlexaTextResponse{}, err
	}

	return NewAlexaTextResponse(fmt.Sprintf("Saved your default bus stop code to %v.", stopCode)), nil
}
