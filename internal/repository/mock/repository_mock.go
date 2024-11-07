// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repositorymock is a generated GoMock package.
package repositorymock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	repository "github.com/iv-sukhanov/finance_tracker/internal/repository"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// AddUsers mocks base method.
func (m *MockUser) AddUsers(users []ftracker.User) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUsers", users)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUsers indicates an expected call of AddUsers.
func (mr *MockUserMockRecorder) AddUsers(users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUsers", reflect.TypeOf((*MockUser)(nil).AddUsers), users)
}

// GetUsers mocks base method.
func (m *MockUser) GetUsers(opts repository.UserOptions) ([]ftracker.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", opts)
	ret0, _ := ret[0].([]ftracker.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserMockRecorder) GetUsers(opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUser)(nil).GetUsers), opts)
}

// MockSpendingCategory is a mock of SpendingCategory interface.
type MockSpendingCategory struct {
	ctrl     *gomock.Controller
	recorder *MockSpendingCategoryMockRecorder
}

// MockSpendingCategoryMockRecorder is the mock recorder for MockSpendingCategory.
type MockSpendingCategoryMockRecorder struct {
	mock *MockSpendingCategory
}

// NewMockSpendingCategory creates a new mock instance.
func NewMockSpendingCategory(ctrl *gomock.Controller) *MockSpendingCategory {
	mock := &MockSpendingCategory{ctrl: ctrl}
	mock.recorder = &MockSpendingCategoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpendingCategory) EXPECT() *MockSpendingCategoryMockRecorder {
	return m.recorder
}

// AddCategories mocks base method.
func (m *MockSpendingCategory) AddCategories(category []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCategories", category)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCategories indicates an expected call of AddCategories.
func (mr *MockSpendingCategoryMockRecorder) AddCategories(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCategories", reflect.TypeOf((*MockSpendingCategory)(nil).AddCategories), category)
}

// GetCategories mocks base method.
func (m *MockSpendingCategory) GetCategories(opts repository.CategoryOptions) ([]ftracker.SpendingCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategories", opts)
	ret0, _ := ret[0].([]ftracker.SpendingCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories.
func (mr *MockSpendingCategoryMockRecorder) GetCategories(opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockSpendingCategory)(nil).GetCategories), opts)
}

// MockSpendingRecord is a mock of SpendingRecord interface.
type MockSpendingRecord struct {
	ctrl     *gomock.Controller
	recorder *MockSpendingRecordMockRecorder
}

// MockSpendingRecordMockRecorder is the mock recorder for MockSpendingRecord.
type MockSpendingRecordMockRecorder struct {
	mock *MockSpendingRecord
}

// NewMockSpendingRecord creates a new mock instance.
func NewMockSpendingRecord(ctrl *gomock.Controller) *MockSpendingRecord {
	mock := &MockSpendingRecord{ctrl: ctrl}
	mock.recorder = &MockSpendingRecordMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpendingRecord) EXPECT() *MockSpendingRecordMockRecorder {
	return m.recorder
}

// AddRecords mocks base method.
func (m *MockSpendingRecord) AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRecords", records)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRecords indicates an expected call of AddRecords.
func (mr *MockSpendingRecordMockRecorder) AddRecords(records interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRecords", reflect.TypeOf((*MockSpendingRecord)(nil).AddRecords), records)
}

// GetRecords mocks base method.
func (m *MockSpendingRecord) GetRecords(opts repository.RecordOptions) ([]ftracker.SpendingRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecords", opts)
	ret0, _ := ret[0].([]ftracker.SpendingRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecords indicates an expected call of GetRecords.
func (mr *MockSpendingRecordMockRecorder) GetRecords(opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecords", reflect.TypeOf((*MockSpendingRecord)(nil).GetRecords), opts)
}
