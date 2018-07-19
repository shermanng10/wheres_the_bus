package main

import (
	"context"
	"log"
)

type Mux interface {
	Handle(context.Context, AlexaRequest) (AlexaTextResponse, error)
}

func AlexaMuxFactory() *AlexaMux {
	return NewAlexaMux(AlexaBusLocationHandlerFactory())
}

func NewAlexaMux(blh BusLocationHandler) *AlexaMux {
	return &AlexaMux{
		busLocationHandler: blh,
	}
}

type AlexaMux struct {
	busLocationHandler BusLocationHandler
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
