// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fn-code/withfun/testing/internal/storage (interfaces: Storage)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	storage "github.com/fn-code/withfun/testing/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockStorage) Add(arg0, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockStorageMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockStorage)(nil).Add), arg0, arg1)
}

// AddCustom mocks base method.
func (m *MockStorage) AddCustom(arg0 *storage.Custom) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCustom", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCustom indicates an expected call of AddCustom.
func (mr *MockStorageMockRecorder) AddCustom(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCustom", reflect.TypeOf((*MockStorage)(nil).AddCustom), arg0)
}