package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
)

func makeTestIntentRequest(intent string) AlexaRequest {
	request := AlexaRequest{}
	request.Version = "1.0"
	request.Request.Type = "IntentRequest"
	request.Request.Time = "2018-07-18T05:42:26Z"
	request.Request.Intent.ConfirmationStatus = "NONE"
	request.Request.Intent.Name = intent
	return request
}

func TestAlexaMuxHandle(t *testing.T) {
	t.Run("Handle wheres_the_bus intent calls bus handler with proper response", func(t *testing.T) {
		log.SetOutput(ioutil.Discard)
		req := makeTestIntentRequest("wheres_the_bus")

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockedResp := NewAlexaTextResponse("There is one bus coming, the Q59 which is 1 stop away.")
		mockBusHandler := NewMockBusLocationHandler(mockCtrl)
		mockBusHandler.EXPECT().GetBusTimes(nil, req).Return(mockedResp, nil).Times(1)

		alexaMux := NewAlexaMux(mockBusHandler)
		resp, _ := alexaMux.Handle(nil, req)
		if resp != mockedResp {
			t.Errorf("Expected response %v, got: %v", mockedResp, resp)
		}
	})
}
