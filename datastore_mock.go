// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package main is a generated GoMock package.
package main

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBusStopPreferenceStore is a mock of BusStopPreferenceStore interface
type MockBusStopPreferenceStore struct {
	ctrl     *gomock.Controller
	recorder *MockBusStopPreferenceStoreMockRecorder
}

// MockBusStopPreferenceStoreMockRecorder is the mock recorder for MockBusStopPreferenceStore
type MockBusStopPreferenceStoreMockRecorder struct {
	mock *MockBusStopPreferenceStore
}

// NewMockBusStopPreferenceStore creates a new mock instance
func NewMockBusStopPreferenceStore(ctrl *gomock.Controller) *MockBusStopPreferenceStore {
	mock := &MockBusStopPreferenceStore{ctrl: ctrl}
	mock.recorder = &MockBusStopPreferenceStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBusStopPreferenceStore) EXPECT() *MockBusStopPreferenceStoreMockRecorder {
	return m.recorder
}

// GetStopCodePreference mocks base method
func (m *MockBusStopPreferenceStore) GetStopCodePreference(userId, prefName string) (string, error) {
	ret := m.ctrl.Call(m, "GetStopCodePreference", userId, prefName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStopCodePreference indicates an expected call of GetStopCodePreference
func (mr *MockBusStopPreferenceStoreMockRecorder) GetStopCodePreference(userId, prefName interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStopCodePreference", reflect.TypeOf((*MockBusStopPreferenceStore)(nil).GetStopCodePreference), userId, prefName)
}

// SetStopCodePreference mocks base method
func (m *MockBusStopPreferenceStore) SetStopCodePreference(userId, prefName, stopCode string) error {
	ret := m.ctrl.Call(m, "SetStopCodePreference", userId, prefName, stopCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStopCodePreference indicates an expected call of SetStopCodePreference
func (mr *MockBusStopPreferenceStoreMockRecorder) SetStopCodePreference(userId, prefName, stopCode interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStopCodePreference", reflect.TypeOf((*MockBusStopPreferenceStore)(nil).SetStopCodePreference), userId, prefName, stopCode)
}
