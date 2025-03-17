// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	service "github.com/iv-sukhanov/finance_tracker/internal/service"
	excelize "github.com/xuri/excelize/v2"
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
func (m *MockUser) GetUsers(opts ...service.UserOption) ([]ftracker.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUsers", varargs...)
	ret0, _ := ret[0].([]ftracker.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserMockRecorder) GetUsers(opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUser)(nil).GetUsers), opts...)
}

// UsersWithGUIDs mocks base method.
func (m *MockUser) UsersWithGUIDs(guids []uuid.UUID) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithGUIDs", guids)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithGUIDs indicates an expected call of UsersWithGUIDs.
func (mr *MockUserMockRecorder) UsersWithGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithGUIDs", reflect.TypeOf((*MockUser)(nil).UsersWithGUIDs), guids)
}

// UsersWithLimit mocks base method.
func (m *MockUser) UsersWithLimit(limit int) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithLimit", limit)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithLimit indicates an expected call of UsersWithLimit.
func (mr *MockUserMockRecorder) UsersWithLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithLimit", reflect.TypeOf((*MockUser)(nil).UsersWithLimit), limit)
}

// UsersWithTelegramIDs mocks base method.
func (m *MockUser) UsersWithTelegramIDs(telegramIDs []string) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithTelegramIDs", telegramIDs)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithTelegramIDs indicates an expected call of UsersWithTelegramIDs.
func (mr *MockUserMockRecorder) UsersWithTelegramIDs(telegramIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithTelegramIDs", reflect.TypeOf((*MockUser)(nil).UsersWithTelegramIDs), telegramIDs)
}

// UsersWithUsernames mocks base method.
func (m *MockUser) UsersWithUsernames(usernames []string) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithUsernames", usernames)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithUsernames indicates an expected call of UsersWithUsernames.
func (mr *MockUserMockRecorder) UsersWithUsernames(usernames interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithUsernames", reflect.TypeOf((*MockUser)(nil).UsersWithUsernames), usernames)
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
func (m *MockSpendingCategory) AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCategories", categories)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCategories indicates an expected call of AddCategories.
func (mr *MockSpendingCategoryMockRecorder) AddCategories(categories interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCategories", reflect.TypeOf((*MockSpendingCategory)(nil).AddCategories), categories)
}

// GetCategories mocks base method.
func (m *MockSpendingCategory) GetCategories(opts ...service.CategoryOption) ([]ftracker.SpendingCategory, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCategories", varargs...)
	ret0, _ := ret[0].([]ftracker.SpendingCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories.
func (mr *MockSpendingCategoryMockRecorder) GetCategories(opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockSpendingCategory)(nil).GetCategories), opts...)
}

// SpendingCategoriesWithCategories mocks base method.
func (m *MockSpendingCategory) SpendingCategoriesWithCategories(categories []string) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithCategories", categories)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithCategories indicates an expected call of SpendingCategoriesWithCategories.
func (mr *MockSpendingCategoryMockRecorder) SpendingCategoriesWithCategories(categories interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithCategories", reflect.TypeOf((*MockSpendingCategory)(nil).SpendingCategoriesWithCategories), categories)
}

// SpendingCategoriesWithGUIDs mocks base method.
func (m *MockSpendingCategory) SpendingCategoriesWithGUIDs(guids []uuid.UUID) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithGUIDs", guids)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithGUIDs indicates an expected call of SpendingCategoriesWithGUIDs.
func (mr *MockSpendingCategoryMockRecorder) SpendingCategoriesWithGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithGUIDs", reflect.TypeOf((*MockSpendingCategory)(nil).SpendingCategoriesWithGUIDs), guids)
}

// SpendingCategoriesWithLimit mocks base method.
func (m *MockSpendingCategory) SpendingCategoriesWithLimit(limit int) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithLimit", limit)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithLimit indicates an expected call of SpendingCategoriesWithLimit.
func (mr *MockSpendingCategoryMockRecorder) SpendingCategoriesWithLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithLimit", reflect.TypeOf((*MockSpendingCategory)(nil).SpendingCategoriesWithLimit), limit)
}

// SpendingCategoriesWithOrder mocks base method.
func (m *MockSpendingCategory) SpendingCategoriesWithOrder(order service.CategoryOrder, asc bool) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithOrder", order, asc)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithOrder indicates an expected call of SpendingCategoriesWithOrder.
func (mr *MockSpendingCategoryMockRecorder) SpendingCategoriesWithOrder(order, asc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithOrder", reflect.TypeOf((*MockSpendingCategory)(nil).SpendingCategoriesWithOrder), order, asc)
}

// SpendingCategoriesWithUserGUIDs mocks base method.
func (m *MockSpendingCategory) SpendingCategoriesWithUserGUIDs(guids []uuid.UUID) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithUserGUIDs", guids)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithUserGUIDs indicates an expected call of SpendingCategoriesWithUserGUIDs.
func (mr *MockSpendingCategoryMockRecorder) SpendingCategoriesWithUserGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithUserGUIDs", reflect.TypeOf((*MockSpendingCategory)(nil).SpendingCategoriesWithUserGUIDs), guids)
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
func (m *MockSpendingRecord) GetRecords(opts ...service.RecordOption) ([]ftracker.SpendingRecord, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRecords", varargs...)
	ret0, _ := ret[0].([]ftracker.SpendingRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecords indicates an expected call of GetRecords.
func (mr *MockSpendingRecordMockRecorder) GetRecords(opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecords", reflect.TypeOf((*MockSpendingRecord)(nil).GetRecords), opts...)
}

// SpendingRecordsWithCategoryGUIDs mocks base method.
func (m *MockSpendingRecord) SpendingRecordsWithCategoryGUIDs(guids []uuid.UUID) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithCategoryGUIDs", guids)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithCategoryGUIDs indicates an expected call of SpendingRecordsWithCategoryGUIDs.
func (mr *MockSpendingRecordMockRecorder) SpendingRecordsWithCategoryGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithCategoryGUIDs", reflect.TypeOf((*MockSpendingRecord)(nil).SpendingRecordsWithCategoryGUIDs), guids)
}

// SpendingRecordsWithGUIDs mocks base method.
func (m *MockSpendingRecord) SpendingRecordsWithGUIDs(guids []uuid.UUID) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithGUIDs", guids)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithGUIDs indicates an expected call of SpendingRecordsWithGUIDs.
func (mr *MockSpendingRecordMockRecorder) SpendingRecordsWithGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithGUIDs", reflect.TypeOf((*MockSpendingRecord)(nil).SpendingRecordsWithGUIDs), guids)
}

// SpendingRecordsWithLimit mocks base method.
func (m *MockSpendingRecord) SpendingRecordsWithLimit(limit int) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithLimit", limit)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithLimit indicates an expected call of SpendingRecordsWithLimit.
func (mr *MockSpendingRecordMockRecorder) SpendingRecordsWithLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithLimit", reflect.TypeOf((*MockSpendingRecord)(nil).SpendingRecordsWithLimit), limit)
}

// SpendingRecordsWithOrder mocks base method.
func (m *MockSpendingRecord) SpendingRecordsWithOrder(order service.RecordOrder, asc bool) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithOrder", order, asc)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithOrder indicates an expected call of SpendingRecordsWithOrder.
func (mr *MockSpendingRecordMockRecorder) SpendingRecordsWithOrder(order, asc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithOrder", reflect.TypeOf((*MockSpendingRecord)(nil).SpendingRecordsWithOrder), order, asc)
}

// SpendingRecordsWithTimeFrame mocks base method.
func (m *MockSpendingRecord) SpendingRecordsWithTimeFrame(from, to time.Time) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithTimeFrame", from, to)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithTimeFrame indicates an expected call of SpendingRecordsWithTimeFrame.
func (mr *MockSpendingRecordMockRecorder) SpendingRecordsWithTimeFrame(from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithTimeFrame", reflect.TypeOf((*MockSpendingRecord)(nil).SpendingRecordsWithTimeFrame), from, to)
}

// MockExelMaker is a mock of ExelMaker interface.
type MockExelMaker struct {
	ctrl     *gomock.Controller
	recorder *MockExelMakerMockRecorder
}

// MockExelMakerMockRecorder is the mock recorder for MockExelMaker.
type MockExelMakerMockRecorder struct {
	mock *MockExelMaker
}

// NewMockExelMaker creates a new mock instance.
func NewMockExelMaker(ctrl *gomock.Controller) *MockExelMaker {
	mock := &MockExelMaker{ctrl: ctrl}
	mock.recorder = &MockExelMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExelMaker) EXPECT() *MockExelMakerMockRecorder {
	return m.recorder
}

// CreateExelFromRecords mocks base method.
func (m *MockExelMaker) CreateExelFromRecords(username string, recods []ftracker.SpendingRecord) (*excelize.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateExelFromRecords", username, recods)
	ret0, _ := ret[0].(*excelize.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateExelFromRecords indicates an expected call of CreateExelFromRecords.
func (mr *MockExelMakerMockRecorder) CreateExelFromRecords(username, recods interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExelFromRecords", reflect.TypeOf((*MockExelMaker)(nil).CreateExelFromRecords), username, recods)
}

// MockServiceInterface is a mock of ServiceInterface interface.
type MockServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockServiceInterfaceMockRecorder
}

// MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.
type MockServiceInterfaceMockRecorder struct {
	mock *MockServiceInterface
}

// NewMockServiceInterface creates a new mock instance.
func NewMockServiceInterface(ctrl *gomock.Controller) *MockServiceInterface {
	mock := &MockServiceInterface{ctrl: ctrl}
	mock.recorder = &MockServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceInterface) EXPECT() *MockServiceInterfaceMockRecorder {
	return m.recorder
}

// AddCategories mocks base method.
func (m *MockServiceInterface) AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCategories", categories)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCategories indicates an expected call of AddCategories.
func (mr *MockServiceInterfaceMockRecorder) AddCategories(categories interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCategories", reflect.TypeOf((*MockServiceInterface)(nil).AddCategories), categories)
}

// AddRecords mocks base method.
func (m *MockServiceInterface) AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRecords", records)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRecords indicates an expected call of AddRecords.
func (mr *MockServiceInterfaceMockRecorder) AddRecords(records interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRecords", reflect.TypeOf((*MockServiceInterface)(nil).AddRecords), records)
}

// AddUsers mocks base method.
func (m *MockServiceInterface) AddUsers(users []ftracker.User) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUsers", users)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUsers indicates an expected call of AddUsers.
func (mr *MockServiceInterfaceMockRecorder) AddUsers(users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUsers", reflect.TypeOf((*MockServiceInterface)(nil).AddUsers), users)
}

// CreateExelFromRecords mocks base method.
func (m *MockServiceInterface) CreateExelFromRecords(username string, recods []ftracker.SpendingRecord) (*excelize.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateExelFromRecords", username, recods)
	ret0, _ := ret[0].(*excelize.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateExelFromRecords indicates an expected call of CreateExelFromRecords.
func (mr *MockServiceInterfaceMockRecorder) CreateExelFromRecords(username, recods interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExelFromRecords", reflect.TypeOf((*MockServiceInterface)(nil).CreateExelFromRecords), username, recods)
}

// GetCategories mocks base method.
func (m *MockServiceInterface) GetCategories(opts ...service.CategoryOption) ([]ftracker.SpendingCategory, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCategories", varargs...)
	ret0, _ := ret[0].([]ftracker.SpendingCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories.
func (mr *MockServiceInterfaceMockRecorder) GetCategories(opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockServiceInterface)(nil).GetCategories), opts...)
}

// GetRecords mocks base method.
func (m *MockServiceInterface) GetRecords(opts ...service.RecordOption) ([]ftracker.SpendingRecord, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRecords", varargs...)
	ret0, _ := ret[0].([]ftracker.SpendingRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecords indicates an expected call of GetRecords.
func (mr *MockServiceInterfaceMockRecorder) GetRecords(opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecords", reflect.TypeOf((*MockServiceInterface)(nil).GetRecords), opts...)
}

// GetUsers mocks base method.
func (m *MockServiceInterface) GetUsers(opts ...service.UserOption) ([]ftracker.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUsers", varargs...)
	ret0, _ := ret[0].([]ftracker.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockServiceInterfaceMockRecorder) GetUsers(opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockServiceInterface)(nil).GetUsers), opts...)
}

// SpendingCategoriesWithCategories mocks base method.
func (m *MockServiceInterface) SpendingCategoriesWithCategories(categories []string) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithCategories", categories)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithCategories indicates an expected call of SpendingCategoriesWithCategories.
func (mr *MockServiceInterfaceMockRecorder) SpendingCategoriesWithCategories(categories interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithCategories", reflect.TypeOf((*MockServiceInterface)(nil).SpendingCategoriesWithCategories), categories)
}

// SpendingCategoriesWithGUIDs mocks base method.
func (m *MockServiceInterface) SpendingCategoriesWithGUIDs(guids []uuid.UUID) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithGUIDs", guids)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithGUIDs indicates an expected call of SpendingCategoriesWithGUIDs.
func (mr *MockServiceInterfaceMockRecorder) SpendingCategoriesWithGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithGUIDs", reflect.TypeOf((*MockServiceInterface)(nil).SpendingCategoriesWithGUIDs), guids)
}

// SpendingCategoriesWithLimit mocks base method.
func (m *MockServiceInterface) SpendingCategoriesWithLimit(limit int) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithLimit", limit)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithLimit indicates an expected call of SpendingCategoriesWithLimit.
func (mr *MockServiceInterfaceMockRecorder) SpendingCategoriesWithLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithLimit", reflect.TypeOf((*MockServiceInterface)(nil).SpendingCategoriesWithLimit), limit)
}

// SpendingCategoriesWithOrder mocks base method.
func (m *MockServiceInterface) SpendingCategoriesWithOrder(order service.CategoryOrder, asc bool) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithOrder", order, asc)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithOrder indicates an expected call of SpendingCategoriesWithOrder.
func (mr *MockServiceInterfaceMockRecorder) SpendingCategoriesWithOrder(order, asc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithOrder", reflect.TypeOf((*MockServiceInterface)(nil).SpendingCategoriesWithOrder), order, asc)
}

// SpendingCategoriesWithUserGUIDs mocks base method.
func (m *MockServiceInterface) SpendingCategoriesWithUserGUIDs(guids []uuid.UUID) service.CategoryOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingCategoriesWithUserGUIDs", guids)
	ret0, _ := ret[0].(service.CategoryOption)
	return ret0
}

// SpendingCategoriesWithUserGUIDs indicates an expected call of SpendingCategoriesWithUserGUIDs.
func (mr *MockServiceInterfaceMockRecorder) SpendingCategoriesWithUserGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingCategoriesWithUserGUIDs", reflect.TypeOf((*MockServiceInterface)(nil).SpendingCategoriesWithUserGUIDs), guids)
}

// SpendingRecordsWithCategoryGUIDs mocks base method.
func (m *MockServiceInterface) SpendingRecordsWithCategoryGUIDs(guids []uuid.UUID) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithCategoryGUIDs", guids)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithCategoryGUIDs indicates an expected call of SpendingRecordsWithCategoryGUIDs.
func (mr *MockServiceInterfaceMockRecorder) SpendingRecordsWithCategoryGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithCategoryGUIDs", reflect.TypeOf((*MockServiceInterface)(nil).SpendingRecordsWithCategoryGUIDs), guids)
}

// SpendingRecordsWithGUIDs mocks base method.
func (m *MockServiceInterface) SpendingRecordsWithGUIDs(guids []uuid.UUID) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithGUIDs", guids)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithGUIDs indicates an expected call of SpendingRecordsWithGUIDs.
func (mr *MockServiceInterfaceMockRecorder) SpendingRecordsWithGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithGUIDs", reflect.TypeOf((*MockServiceInterface)(nil).SpendingRecordsWithGUIDs), guids)
}

// SpendingRecordsWithLimit mocks base method.
func (m *MockServiceInterface) SpendingRecordsWithLimit(limit int) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithLimit", limit)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithLimit indicates an expected call of SpendingRecordsWithLimit.
func (mr *MockServiceInterfaceMockRecorder) SpendingRecordsWithLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithLimit", reflect.TypeOf((*MockServiceInterface)(nil).SpendingRecordsWithLimit), limit)
}

// SpendingRecordsWithOrder mocks base method.
func (m *MockServiceInterface) SpendingRecordsWithOrder(order service.RecordOrder, asc bool) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithOrder", order, asc)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithOrder indicates an expected call of SpendingRecordsWithOrder.
func (mr *MockServiceInterfaceMockRecorder) SpendingRecordsWithOrder(order, asc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithOrder", reflect.TypeOf((*MockServiceInterface)(nil).SpendingRecordsWithOrder), order, asc)
}

// SpendingRecordsWithTimeFrame mocks base method.
func (m *MockServiceInterface) SpendingRecordsWithTimeFrame(from, to time.Time) service.RecordOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendingRecordsWithTimeFrame", from, to)
	ret0, _ := ret[0].(service.RecordOption)
	return ret0
}

// SpendingRecordsWithTimeFrame indicates an expected call of SpendingRecordsWithTimeFrame.
func (mr *MockServiceInterfaceMockRecorder) SpendingRecordsWithTimeFrame(from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendingRecordsWithTimeFrame", reflect.TypeOf((*MockServiceInterface)(nil).SpendingRecordsWithTimeFrame), from, to)
}

// UsersWithGUIDs mocks base method.
func (m *MockServiceInterface) UsersWithGUIDs(guids []uuid.UUID) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithGUIDs", guids)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithGUIDs indicates an expected call of UsersWithGUIDs.
func (mr *MockServiceInterfaceMockRecorder) UsersWithGUIDs(guids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithGUIDs", reflect.TypeOf((*MockServiceInterface)(nil).UsersWithGUIDs), guids)
}

// UsersWithLimit mocks base method.
func (m *MockServiceInterface) UsersWithLimit(limit int) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithLimit", limit)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithLimit indicates an expected call of UsersWithLimit.
func (mr *MockServiceInterfaceMockRecorder) UsersWithLimit(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithLimit", reflect.TypeOf((*MockServiceInterface)(nil).UsersWithLimit), limit)
}

// UsersWithTelegramIDs mocks base method.
func (m *MockServiceInterface) UsersWithTelegramIDs(telegramIDs []string) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithTelegramIDs", telegramIDs)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithTelegramIDs indicates an expected call of UsersWithTelegramIDs.
func (mr *MockServiceInterfaceMockRecorder) UsersWithTelegramIDs(telegramIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithTelegramIDs", reflect.TypeOf((*MockServiceInterface)(nil).UsersWithTelegramIDs), telegramIDs)
}

// UsersWithUsernames mocks base method.
func (m *MockServiceInterface) UsersWithUsernames(usernames []string) service.UserOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersWithUsernames", usernames)
	ret0, _ := ret[0].(service.UserOption)
	return ret0
}

// UsersWithUsernames indicates an expected call of UsersWithUsernames.
func (mr *MockServiceInterfaceMockRecorder) UsersWithUsernames(usernames interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersWithUsernames", reflect.TypeOf((*MockServiceInterface)(nil).UsersWithUsernames), usernames)
}
