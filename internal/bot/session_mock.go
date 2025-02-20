// Code generated by MockGen. DO NOT EDIT.
// Source: session.go

// Package mock_bot is a generated GoMock package.
package bot

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessions is a mock of Sessions interface.
type MockSessions struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsMockRecorder
}

// MockSessionsMockRecorder is the mock recorder for MockSessions.
type MockSessionsMockRecorder struct {
	mock *MockSessions
}

// NewMockSessions creates a new mock instance.
func NewMockSessions(ctrl *gomock.Controller) *MockSessions {
	mock := &MockSessions{ctrl: ctrl}
	mock.recorder = &MockSessionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessions) EXPECT() *MockSessionsMockRecorder {
	return m.recorder
}

// AddSession mocks base method.
func (m *MockSessions) AddSession(chatID, userID int64, username string) *Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSession", chatID, userID, username)
	ret0, _ := ret[0].(*Session)
	return ret0
}

// AddSession indicates an expected call of AddSession.
func (mr *MockSessionsMockRecorder) AddSession(chatID, userID, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSession", reflect.TypeOf((*MockSessions)(nil).AddSession), chatID, userID, username)
}

// GetSession mocks base method.
func (m *MockSessions) GetSession(chatID int64) *Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", chatID)
	ret0, _ := ret[0].(*Session)
	return ret0
}

// GetSession indicates an expected call of GetSession.
func (mr *MockSessionsMockRecorder) GetSession(chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockSessions)(nil).GetSession), chatID)
}
