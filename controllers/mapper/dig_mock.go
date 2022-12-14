// Code generated by MockGen. DO NOT EDIT.
// Source: controllers/utils/dns.go

// Package mapper is a generated GoMock package.
package mapper

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDigger is a mock of Digger interface.
type MockDigger struct {
	ctrl     *gomock.Controller
	recorder *MockDiggerMockRecorder
}

// MockDiggerMockRecorder is the mock recorder for MockDigger.
type MockDiggerMockRecorder struct {
	mock *MockDigger
}

// NewMockDigger creates a new mock instance.
func NewMockDigger(ctrl *gomock.Controller) *MockDigger {
	mock := &MockDigger{ctrl: ctrl}
	mock.recorder = &MockDiggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDigger) EXPECT() *MockDiggerMockRecorder {
	return m.recorder
}

// DigA mocks base method.
func (m *MockDigger) DigA(fqdn string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DigA", fqdn)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DigA indicates an expected call of DigA.
func (mr *MockDiggerMockRecorder) DigA(fqdn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DigA", reflect.TypeOf((*MockDigger)(nil).DigA), fqdn)
}
