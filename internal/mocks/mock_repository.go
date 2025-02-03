// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/RomanAgaltsev/urlcut/internal/interfaces (interfaces: Repository)
//
// Generated by this command:
//
//	mockgen -destination=internal/mocks/mock_repository.go -package=mocks github.com/RomanAgaltsev/urlcut/internal/interfaces Repository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"

	model "github.com/RomanAgaltsev/urlcut/internal/model"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
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

// Close mocks base method.
func (m *MockRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}

// DeleteURLs mocks base method.
func (m *MockRepository) DeleteURLs(ctx context.Context, urls []*model.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURLs", ctx, urls)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteURLs indicates an expected call of DeleteURLs.
func (mr *MockRepositoryMockRecorder) DeleteURLs(ctx, urls any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURLs", reflect.TypeOf((*MockRepository)(nil).DeleteURLs), ctx, urls)
}

// Get mocks base method.
func (m *MockRepository) Get(ctx context.Context, id string) (*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), ctx, id)
}

// GetUserURLs mocks base method.
func (m *MockRepository) GetUserURLs(ctx context.Context, uid uuid.UUID) ([]*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", ctx, uid)
	ret0, _ := ret[0].([]*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockRepositoryMockRecorder) GetUserURLs(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockRepository)(nil).GetUserURLs), ctx, uid)
}

// Store mocks base method.
func (m *MockRepository) Store(ctx context.Context, urls []*model.URL) (*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, urls)
	ret0, _ := ret[0].(*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockRepositoryMockRecorder) Store(ctx, urls any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockRepository)(nil).Store), ctx, urls)
}
