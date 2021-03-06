// Code generated by MockGen. DO NOT EDIT.
// Source: alexa_bus_location_handler.go

// Package main is a generated GoMock package.
package main

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBusLocationHandler is a mock of BusLocationHandler interface
type MockBusLocationHandler struct {
	ctrl     *gomock.Controller
	recorder *MockBusLocationHandlerMockRecorder
}

// MockBusLocationHandlerMockRecorder is the mock recorder for MockBusLocationHandler
type MockBusLocationHandlerMockRecorder struct {
	mock *MockBusLocationHandler
}

// NewMockBusLocationHandler creates a new mock instance
func NewMockBusLocationHandler(ctrl *gomock.Controller) *MockBusLocationHandler {
	mock := &MockBusLocationHandler{ctrl: ctrl}
	mock.recorder = &MockBusLocationHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBusLocationHandler) EXPECT() *MockBusLocationHandlerMockRecorder {
	return m.recorder
}

// GetBusTimes mocks base method
func (m *MockBusLocationHandler) GetBusTimes(arg0 context.Context, arg1 AlexaRequest) (AlexaTextResponse, error) {
	ret := m.ctrl.Call(m, "GetBusTimes", arg0, arg1)
	ret0, _ := ret[0].(AlexaTextResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBusTimes indicates an expected call of GetBusTimes
func (mr *MockBusLocationHandlerMockRecorder) GetBusTimes(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBusTimes", reflect.TypeOf((*MockBusLocationHandler)(nil).GetBusTimes), arg0, arg1)
}
