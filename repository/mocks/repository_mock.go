package mocks

import (
	models "MartellX/avito-tech-task/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRepository) Delete(arg0 *models.Offer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", arg0)
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), arg0)
}

// FindOffer mocks base method.
func (m *MockRepository) FindOffer(arg0, arg1 uint64) (*models.Offer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOffer", arg0, arg1)
	ret0, _ := ret[0].(*models.Offer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOffer indicates an expected call of FindOffer.
func (mr *MockRepositoryMockRecorder) FindOffer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOffer", reflect.TypeOf((*MockRepository)(nil).FindOffer), arg0, arg1)
}

// FindOffersByConditions mocks base method.
func (m *MockRepository) FindOffersByConditions(arg0 map[string]interface{}) ([]models.Offer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOffersByConditions", arg0)
	ret0, _ := ret[0].([]models.Offer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOffersByConditions indicates an expected call of FindOffersByConditions.
func (mr *MockRepositoryMockRecorder) FindOffersByConditions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOffersByConditions", reflect.TypeOf((*MockRepository)(nil).FindOffersByConditions), arg0)
}

// GetDB mocks base method.
func (m *MockRepository) GetDB() *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDB")
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// GetDB indicates an expected call of GetDB.
func (mr *MockRepositoryMockRecorder) GetDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDB", reflect.TypeOf((*MockRepository)(nil).GetDB))
}

// NewOffer mocks base method.
func (m *MockRepository) NewOffer(arg0, arg1 uint64, arg2 string, arg3 int64, arg4 int, arg5 bool) (*models.Offer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewOffer", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*models.Offer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewOffer indicates an expected call of NewOffer.
func (mr *MockRepositoryMockRecorder) NewOffer(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewOffer", reflect.TypeOf((*MockRepository)(nil).NewOffer), arg0, arg1, arg2, arg3, arg4, arg5)
}

// SetDB mocks base method.
func (m *MockRepository) SetDB(arg0 *gorm.DB) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDB", arg0)
}

// SetDB indicates an expected call of SetDB.
func (mr *MockRepositoryMockRecorder) SetDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDB", reflect.TypeOf((*MockRepository)(nil).SetDB), arg0)
}

// Update mocks base method.
func (m *MockRepository) Update(arg0 *models.Offer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), arg0)
}

// UpdateColumns mocks base method.
func (m *MockRepository) UpdateColumns(arg0 *models.Offer, arg1 string, arg2 int64, arg3 int, arg4 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateColumns", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateColumns indicates an expected call of UpdateColumns.
func (mr *MockRepositoryMockRecorder) UpdateColumns(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateColumns", reflect.TypeOf((*MockRepository)(nil).UpdateColumns), arg0, arg1, arg2, arg3, arg4)
}
