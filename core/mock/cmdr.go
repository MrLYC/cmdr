// Code generated by MockGen. DO NOT EDIT.
// Source: cmdr.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	core "github.com/mrlyc/cmdr/core"
)

// MockCmdrSearcher is a mock of CmdrSearcher interface.
type MockCmdrSearcher struct {
	ctrl     *gomock.Controller
	recorder *MockCmdrSearcherMockRecorder
}

// MockCmdrSearcherMockRecorder is the mock recorder for MockCmdrSearcher.
type MockCmdrSearcherMockRecorder struct {
	mock *MockCmdrSearcher
}

// NewMockCmdrSearcher creates a new mock instance.
func NewMockCmdrSearcher(ctrl *gomock.Controller) *MockCmdrSearcher {
	mock := &MockCmdrSearcher{ctrl: ctrl}
	mock.recorder = &MockCmdrSearcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCmdrSearcher) EXPECT() *MockCmdrSearcherMockRecorder {
	return m.recorder
}

// GetReleaseAsset mocks base method.
func (m *MockCmdrSearcher) GetReleaseAsset(ctx context.Context, releaseName, assetName string) (core.CmdrReleaseAsset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleaseAsset", ctx, releaseName, assetName)
	ret0, _ := ret[0].(core.CmdrReleaseAsset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReleaseAsset indicates an expected call of GetReleaseAsset.
func (mr *MockCmdrSearcherMockRecorder) GetReleaseAsset(ctx, releaseName, assetName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleaseAsset", reflect.TypeOf((*MockCmdrSearcher)(nil).GetReleaseAsset), ctx, releaseName, assetName)
}
