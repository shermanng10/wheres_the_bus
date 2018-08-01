package main

import (
	"context"
	"log"
)

type Mux interface {
	Handle(context.Context, AlexaRequest) (AlexaTextResponse, error)
}

func InitAlexaMux() *AlexaMux {
	return NewAlexaMux(InitAlexaBusLocationHandler(), InitAlexaSetPreferenceHandler())
}

func NewAlexaMux(blh BusLocationHandler, sph SetPreferenceHandler) *AlexaMux {
	return &AlexaMux{
		busLocationHandler:   blh,
		setPreferenceHandler: sph,
	}
}

type AlexaMux struct {
	busLocationHandler   BusLocationHandler
	setPreferenceHandler SetPreferenceHandler
}

func (a *AlexaMux) Handle(ctx context.Context, r AlexaRequest) (AlexaTextResponse, error) {
	log.Println("-------- Alexa Request Parameters: --------")
	log.Printf("AlexaRequest: %+v\n", r)
	log.Println("-------- Done. --------")

	var resp AlexaTextResponse
	var err error

	switch r.Request.Intent.Name {
	case "wheres_the_bus":
		resp, err = a.busLocationHandler.GetBusTimes(ctx, r)
	case "save_bus_code":
		resp, err = a.setPreferenceHandler.SetBusPreference(ctx, r)
	case "AMAZON.HelpIntent":
		resp = NewAlexaTextResponse("Just ask where the bus is.")
	default:
		resp = NewAlexaTextResponse("I'm sorry, I didn't understand that.")
	}

	if err != nil {
		log.Printf("Got Error: %v", err)
		resp = NewAlexaTextResponse("Sorry, something went wrong. Forgive me human.")
	}

	return resp, nil
}
