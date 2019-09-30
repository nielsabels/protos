// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/manager.go

// Package app is a generated GoMock package.
package app

import (
	gomock "github.com/golang/mock/gomock"
	core "protos/internal/core"
	reflect "reflect"
)

// MockappStore is a mock of appStore interface
type MockappStore struct {
	ctrl     *gomock.Controller
	recorder *MockappStoreMockRecorder
}

// MockappStoreMockRecorder is the mock recorder for MockappStore
type MockappStoreMockRecorder struct {
	mock *MockappStore
}

// NewMockappStore creates a new mock instance
func NewMockappStore(ctrl *gomock.Controller) *MockappStore {
	mock := &MockappStore{ctrl: ctrl}
	mock.recorder = &MockappStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockappStore) EXPECT() *MockappStoreMockRecorder {
	return m.recorder
}

// GetInstaller mocks base method
func (m *MockappStore) GetInstaller(arg0 string) (core.Installer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstaller", arg0)
	ret0, _ := ret[0].(core.Installer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstaller indicates an expected call of GetInstaller
func (mr *MockappStoreMockRecorder) GetInstaller(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstaller", reflect.TypeOf((*MockappStore)(nil).GetInstaller), arg0)
}

// MockdnsResource is a mock of dnsResource interface
type MockdnsResource struct {
	ctrl     *gomock.Controller
	recorder *MockdnsResourceMockRecorder
}

// MockdnsResourceMockRecorder is the mock recorder for MockdnsResource
type MockdnsResourceMockRecorder struct {
	mock *MockdnsResource
}

// NewMockdnsResource creates a new mock instance
func NewMockdnsResource(ctrl *gomock.Controller) *MockdnsResource {
	mock := &MockdnsResource{ctrl: ctrl}
	mock.recorder = &MockdnsResourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockdnsResource) EXPECT() *MockdnsResourceMockRecorder {
	return m.recorder
}

// GetName mocks base method
func (m *MockdnsResource) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockdnsResourceMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockdnsResource)(nil).GetName))
}

// GetValue mocks base method
func (m *MockdnsResource) GetValue() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetValue indicates an expected call of GetValue
func (mr *MockdnsResourceMockRecorder) GetValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockdnsResource)(nil).GetValue))
}

// Update mocks base method
func (m *MockdnsResource) Update(arg0 core.Type) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update
func (mr *MockdnsResourceMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockdnsResource)(nil).Update), arg0)
}

// Sanitize mocks base method
func (m *MockdnsResource) Sanitize() core.Type {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sanitize")
	ret0, _ := ret[0].(core.Type)
	return ret0
}

// Sanitize indicates an expected call of Sanitize
func (mr *MockdnsResourceMockRecorder) Sanitize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sanitize", reflect.TypeOf((*MockdnsResource)(nil).Sanitize))
}
