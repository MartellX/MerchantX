// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	services "MartellX/avito-tech-task/services"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskService is a mock of TaskService interface.
type MockTaskService struct {
	ctrl     *gomock.Controller
	recorder *MockTaskServiceMockRecorder
}

// MockTaskServiceMockRecorder is the mock recorder for MockTaskService.
type MockTaskServiceMockRecorder struct {
	mock *MockTaskService
}

// NewMockTaskService creates a new mock instance.
func NewMockTaskService(ctrl *gomock.Controller) *MockTaskService {
	mock := &MockTaskService{ctrl: ctrl}
	mock.recorder = &MockTaskServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskService) EXPECT() *MockTaskServiceMockRecorder {
	return m.recorder
}

// GetTask mock_services base method.
func (m *MockTaskService) GetTask(id string) (*services.Task, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", id)
	ret0, _ := ret[0].(*services.Task)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskServiceMockRecorder) GetTask(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskService)(nil).GetTask), id)
}

// StartUploadingTask mock_services base method.
func (m *MockTaskService) StartUploadingTask(sellerId uint64, xlsxURL string) (*services.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartUploadingTask", sellerId, xlsxURL)
	ret0, _ := ret[0].(*services.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartUploadingTask indicates an expected call of StartUploadingTask.
func (mr *MockTaskServiceMockRecorder) StartUploadingTask(sellerId, xlsxURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartUploadingTask", reflect.TypeOf((*MockTaskService)(nil).StartUploadingTask), sellerId, xlsxURL)
}
