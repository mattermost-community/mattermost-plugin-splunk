// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bakurits/mattermost-plugin-splunk/server/store (interfaces: Store)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	store "github.com/bakurits/mattermost-plugin-splunk/server/store"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// ChangeCurrentUser mocks base method.
func (m *MockStore) ChangeCurrentUser(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeCurrentUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeCurrentUser indicates an expected call of ChangeCurrentUser.
func (mr *MockStoreMockRecorder) ChangeCurrentUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeCurrentUser", reflect.TypeOf((*MockStore)(nil).ChangeCurrentUser), arg0, arg1)
}

// CurrentUser mocks base method.
func (m *MockStore) CurrentUser(arg0 string) (store.SplunkUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentUser", arg0)
	ret0, _ := ret[0].(store.SplunkUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrentUser indicates an expected call of CurrentUser.
func (mr *MockStoreMockRecorder) CurrentUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentUser", reflect.TypeOf((*MockStore)(nil).CurrentUser), arg0)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1, arg2)
}

// RegisterUser mocks base method.
func (m *MockStore) RegisterUser(arg0 string, arg1 store.SplunkUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockStoreMockRecorder) RegisterUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockStore)(nil).RegisterUser), arg0, arg1)
}

// User mocks base method.
func (m *MockStore) User(arg0, arg1, arg2 string) (store.SplunkUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User", arg0, arg1, arg2)
	ret0, _ := ret[0].(store.SplunkUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// User indicates an expected call of User.
func (mr *MockStoreMockRecorder) User(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockStore)(nil).User), arg0, arg1, arg2)
}