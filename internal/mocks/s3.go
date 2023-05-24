// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/T-V-N/gophkeeper/internal/app (interfaces: S3Store)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockS3Store is a mock of S3Store interface.
type MockS3Store struct {
	ctrl     *gomock.Controller
	recorder *MockS3StoreMockRecorder
}

// MockS3StoreMockRecorder is the mock recorder for MockS3Store.
type MockS3StoreMockRecorder struct {
	mock *MockS3Store
}

// NewMockS3Store creates a new mock instance.
func NewMockS3Store(ctrl *gomock.Controller) *MockS3Store {
	mock := &MockS3Store{ctrl: ctrl}
	mock.recorder = &MockS3StoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockS3Store) EXPECT() *MockS3StoreMockRecorder {
	return m.recorder
}

// GetFileInfo mocks base method.
func (m *MockS3Store) GetFileInfo(arg0 context.Context, arg1 string) (time.Time, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileInfo", arg0, arg1)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFileInfo indicates an expected call of GetFileInfo.
func (mr *MockS3StoreMockRecorder) GetFileInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileInfo", reflect.TypeOf((*MockS3Store)(nil).GetFileInfo), arg0, arg1)
}

// GetUploadLink mocks base method.
func (m *MockS3Store) GetUploadLink(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUploadLink", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUploadLink indicates an expected call of GetUploadLink.
func (mr *MockS3StoreMockRecorder) GetUploadLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUploadLink", reflect.TypeOf((*MockS3Store)(nil).GetUploadLink), arg0, arg1)
}
