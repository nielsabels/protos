// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/capability.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	core "github.com/protosio/protos/internal/core"
	reflect "reflect"
)

// MockCapabilityManager is a mock of CapabilityManager interface
type MockCapabilityManager struct {
	ctrl     *gomock.Controller
	recorder *MockCapabilityManagerMockRecorder
}

// MockCapabilityManagerMockRecorder is the mock recorder for MockCapabilityManager
type MockCapabilityManagerMockRecorder struct {
	mock *MockCapabilityManager
}

// NewMockCapabilityManager creates a new mock instance
func NewMockCapabilityManager(ctrl *gomock.Controller) *MockCapabilityManager {
	mock := &MockCapabilityManager{ctrl: ctrl}
	mock.recorder = &MockCapabilityManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCapabilityManager) EXPECT() *MockCapabilityManagerMockRecorder {
	return m.recorder
}

// Validate mocks base method
func (m *MockCapabilityManager) Validate(methodcap core.Capability, appcap string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", methodcap, appcap)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockCapabilityManagerMockRecorder) Validate(methodcap, appcap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockCapabilityManager)(nil).Validate), methodcap, appcap)
}

// SetMethodCap mocks base method
func (m *MockCapabilityManager) SetMethodCap(method string, cap core.Capability) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMethodCap", method, cap)
}

// SetMethodCap indicates an expected call of SetMethodCap
func (mr *MockCapabilityManagerMockRecorder) SetMethodCap(method, cap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMethodCap", reflect.TypeOf((*MockCapabilityManager)(nil).SetMethodCap), method, cap)
}

// GetMethodCap mocks base method
func (m *MockCapabilityManager) GetMethodCap(method string) (core.Capability, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMethodCap", method)
	ret0, _ := ret[0].(core.Capability)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMethodCap indicates an expected call of GetMethodCap
func (mr *MockCapabilityManagerMockRecorder) GetMethodCap(method interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMethodCap", reflect.TypeOf((*MockCapabilityManager)(nil).GetMethodCap), method)
}

// GetByName mocks base method
func (m *MockCapabilityManager) GetByName(name string) (core.Capability, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", name)
	ret0, _ := ret[0].(core.Capability)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName
func (mr *MockCapabilityManagerMockRecorder) GetByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockCapabilityManager)(nil).GetByName), name)
}

// GetOrPanic mocks base method
func (m *MockCapabilityManager) GetOrPanic(name string) core.Capability {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrPanic", name)
	ret0, _ := ret[0].(core.Capability)
	return ret0
}

// GetOrPanic indicates an expected call of GetOrPanic
func (mr *MockCapabilityManagerMockRecorder) GetOrPanic(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrPanic", reflect.TypeOf((*MockCapabilityManager)(nil).GetOrPanic), name)
}

// ClearAll mocks base method
func (m *MockCapabilityManager) ClearAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearAll")
}

// ClearAll indicates an expected call of ClearAll
func (mr *MockCapabilityManagerMockRecorder) ClearAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearAll", reflect.TypeOf((*MockCapabilityManager)(nil).ClearAll))
}

// MockCapability is a mock of Capability interface
type MockCapability struct {
	ctrl     *gomock.Controller
	recorder *MockCapabilityMockRecorder
}

// MockCapabilityMockRecorder is the mock recorder for MockCapability
type MockCapabilityMockRecorder struct {
	mock *MockCapability
}

// NewMockCapability creates a new mock instance
func NewMockCapability(ctrl *gomock.Controller) *MockCapability {
	mock := &MockCapability{ctrl: ctrl}
	mock.recorder = &MockCapabilityMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCapability) EXPECT() *MockCapabilityMockRecorder {
	return m.recorder
}

// GetName mocks base method
func (m *MockCapability) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockCapabilityMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockCapability)(nil).GetName))
}

// GetParent mocks base method
func (m *MockCapability) GetParent() core.Capability {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParent")
	ret0, _ := ret[0].(core.Capability)
	return ret0
}

// GetParent indicates an expected call of GetParent
func (mr *MockCapabilityMockRecorder) GetParent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParent", reflect.TypeOf((*MockCapability)(nil).GetParent))
}

// MockCapabilityChecker is a mock of CapabilityChecker interface
type MockCapabilityChecker struct {
	ctrl     *gomock.Controller
	recorder *MockCapabilityCheckerMockRecorder
}

// MockCapabilityCheckerMockRecorder is the mock recorder for MockCapabilityChecker
type MockCapabilityCheckerMockRecorder struct {
	mock *MockCapabilityChecker
}

// NewMockCapabilityChecker creates a new mock instance
func NewMockCapabilityChecker(ctrl *gomock.Controller) *MockCapabilityChecker {
	mock := &MockCapabilityChecker{ctrl: ctrl}
	mock.recorder = &MockCapabilityCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCapabilityChecker) EXPECT() *MockCapabilityCheckerMockRecorder {
	return m.recorder
}

// ValidateCapability mocks base method
func (m *MockCapabilityChecker) ValidateCapability(cap core.Capability) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCapability", cap)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateCapability indicates an expected call of ValidateCapability
func (mr *MockCapabilityCheckerMockRecorder) ValidateCapability(cap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCapability", reflect.TypeOf((*MockCapabilityChecker)(nil).ValidateCapability), cap)
}
