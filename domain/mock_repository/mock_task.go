// Code generated by MockGen. DO NOT EDIT.
// Source: task.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	gomock "github.com/golang/mock/gomock"
	entity "github.com/hiroyaonoe/todoapp-server/domain/entity"
	gorm "github.com/jinzhu/gorm"
	reflect "reflect"
)

// MockTaskRepository is a mock of TaskRepository interface
type MockTaskRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTaskRepositoryMockRecorder
}

// MockTaskRepositoryMockRecorder is the mock recorder for MockTaskRepository
type MockTaskRepositoryMockRecorder struct {
	mock *MockTaskRepository
}

// NewMockTaskRepository creates a new mock instance
func NewMockTaskRepository(ctrl *gomock.Controller) *MockTaskRepository {
	mock := &MockTaskRepository{ctrl: ctrl}
	mock.recorder = &MockTaskRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskRepository) EXPECT() *MockTaskRepositoryMockRecorder {
	return m.recorder
}

// FindByUser mocks base method
func (m *MockTaskRepository) FindByUser(db *gorm.DB, uid int) ([]entity.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUser", db, uid)
	ret0, _ := ret[0].([]entity.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUser indicates an expected call of FindByUser
func (mr *MockTaskRepositoryMockRecorder) FindByUser(db, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUser", reflect.TypeOf((*MockTaskRepository)(nil).FindByUser), db, uid)
}

// FindByID mocks base method
func (m *MockTaskRepository) FindByID(db *gorm.DB, id int) (entity.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", db, id)
	ret0, _ := ret[0].(entity.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockTaskRepositoryMockRecorder) FindByID(db, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockTaskRepository)(nil).FindByID), db, id)
}

// Create mocks base method
func (m *MockTaskRepository) Create(db *gorm.DB, t entity.Task) (entity.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", db, t)
	ret0, _ := ret[0].(entity.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockTaskRepositoryMockRecorder) Create(db, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTaskRepository)(nil).Create), db, t)
}

// Update mocks base method
func (m *MockTaskRepository) Update(db *gorm.DB, t entity.Task) (entity.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", db, t)
	ret0, _ := ret[0].(entity.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockTaskRepositoryMockRecorder) Update(db, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTaskRepository)(nil).Update), db, t)
}

// Delete mocks base method
func (m *MockTaskRepository) Delete(db *gorm.DB, id int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", db, id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete
func (mr *MockTaskRepositoryMockRecorder) Delete(db, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTaskRepository)(nil).Delete), db, id)
}