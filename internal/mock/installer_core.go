// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/installer.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	core "github.com/protosio/protos/internal/core"
	reflect "reflect"
)

// MockAppStore is a mock of AppStore interface
type MockAppStore struct {
	ctrl     *gomock.Controller
	recorder *MockAppStoreMockRecorder
}

// MockAppStoreMockRecorder is the mock recorder for MockAppStore
type MockAppStoreMockRecorder struct {
	mock *MockAppStore
}

// NewMockAppStore creates a new mock instance
func NewMockAppStore(ctrl *gomock.Controller) *MockAppStore {
	mock := &MockAppStore{ctrl: ctrl}
	mock.recorder = &MockAppStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppStore) EXPECT() *MockAppStoreMockRecorder {
	return m.recorder
}

// GetInstallers mocks base method
func (m *MockAppStore) GetInstallers() (map[string]core.Installer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstallers")
	ret0, _ := ret[0].(map[string]core.Installer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstallers indicates an expected call of GetInstallers
func (mr *MockAppStoreMockRecorder) GetInstallers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstallers", reflect.TypeOf((*MockAppStore)(nil).GetInstallers))
}

// GetInstaller mocks base method
func (m *MockAppStore) GetInstaller(id string) (core.Installer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstaller", id)
	ret0, _ := ret[0].(core.Installer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstaller indicates an expected call of GetInstaller
func (mr *MockAppStoreMockRecorder) GetInstaller(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstaller", reflect.TypeOf((*MockAppStore)(nil).GetInstaller), id)
}

// Search mocks base method
func (m *MockAppStore) Search(key, value string) (map[string]core.Installer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", key, value)
	ret0, _ := ret[0].(map[string]core.Installer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockAppStoreMockRecorder) Search(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockAppStore)(nil).Search), key, value)
}

// MockInstallerCache is a mock of InstallerCache interface
type MockInstallerCache struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerCacheMockRecorder
}

// MockInstallerCacheMockRecorder is the mock recorder for MockInstallerCache
type MockInstallerCacheMockRecorder struct {
	mock *MockInstallerCache
}

// NewMockInstallerCache creates a new mock instance
func NewMockInstallerCache(ctrl *gomock.Controller) *MockInstallerCache {
	mock := &MockInstallerCache{ctrl: ctrl}
	mock.recorder = &MockInstallerCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInstallerCache) EXPECT() *MockInstallerCacheMockRecorder {
	return m.recorder
}

// GetLocalInstallers mocks base method
func (m *MockInstallerCache) GetLocalInstallers() (map[string]core.Installer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalInstallers")
	ret0, _ := ret[0].(map[string]core.Installer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocalInstallers indicates an expected call of GetLocalInstallers
func (mr *MockInstallerCacheMockRecorder) GetLocalInstallers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalInstallers", reflect.TypeOf((*MockInstallerCache)(nil).GetLocalInstallers))
}

// GetLocalInstaller mocks base method
func (m *MockInstallerCache) GetLocalInstaller(id string) (core.Installer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalInstaller", id)
	ret0, _ := ret[0].(core.Installer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocalInstaller indicates an expected call of GetLocalInstaller
func (mr *MockInstallerCacheMockRecorder) GetLocalInstaller(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalInstaller", reflect.TypeOf((*MockInstallerCache)(nil).GetLocalInstaller), id)
}

// RemoveLocalInstaller mocks base method
func (m *MockInstallerCache) RemoveLocalInstaller(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLocalInstaller", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLocalInstaller indicates an expected call of RemoveLocalInstaller
func (mr *MockInstallerCacheMockRecorder) RemoveLocalInstaller(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLocalInstaller", reflect.TypeOf((*MockInstallerCache)(nil).RemoveLocalInstaller), id)
}

// MockInstaller is a mock of Installer interface
type MockInstaller struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerMockRecorder
}

// MockInstallerMockRecorder is the mock recorder for MockInstaller
type MockInstallerMockRecorder struct {
	mock *MockInstaller
}

// NewMockInstaller creates a new mock instance
func NewMockInstaller(ctrl *gomock.Controller) *MockInstaller {
	mock := &MockInstaller{ctrl: ctrl}
	mock.recorder = &MockInstallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInstaller) EXPECT() *MockInstallerMockRecorder {
	return m.recorder
}

// GetLastVersion mocks base method
func (m *MockInstaller) GetLastVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLastVersion indicates an expected call of GetLastVersion
func (mr *MockInstallerMockRecorder) GetLastVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastVersion", reflect.TypeOf((*MockInstaller)(nil).GetLastVersion))
}

// GetMetadata mocks base method
func (m *MockInstaller) GetMetadata(version string) (core.InstallerMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadata", version)
	ret0, _ := ret[0].(core.InstallerMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata
func (mr *MockInstallerMockRecorder) GetMetadata(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockInstaller)(nil).GetMetadata), version)
}

// IsPlatformImageAvailable mocks base method
func (m *MockInstaller) IsPlatformImageAvailable(version string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPlatformImageAvailable", version)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPlatformImageAvailable indicates an expected call of IsPlatformImageAvailable
func (mr *MockInstallerMockRecorder) IsPlatformImageAvailable(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPlatformImageAvailable", reflect.TypeOf((*MockInstaller)(nil).IsPlatformImageAvailable), version)
}

// DownloadAsync mocks base method
func (m *MockInstaller) DownloadAsync(version, appID string) core.Task {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadAsync", version, appID)
	ret0, _ := ret[0].(core.Task)
	return ret0
}

// DownloadAsync indicates an expected call of DownloadAsync
func (mr *MockInstallerMockRecorder) DownloadAsync(version, appID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadAsync", reflect.TypeOf((*MockInstaller)(nil).DownloadAsync), version, appID)
}
