// Code generated by MockGen. DO NOT EDIT.
// Source: root.go
//
// Generated by this command:
//
//	mockgen -source=root.go -destination=mock_root.go -package=cmd
//

// Package cmd is a generated GoMock package.
package cmd

import (
	reflect "reflect"

	query "github.com/go-zen-chu/aictl/usecase/query"
	gomock "go.uber.org/mock/gomock"
)

// MockCommandRequirements is a mock of CommandRequirements interface.
type MockCommandRequirements struct {
	ctrl     *gomock.Controller
	recorder *MockCommandRequirementsMockRecorder
}

// MockCommandRequirementsMockRecorder is the mock recorder for MockCommandRequirements.
type MockCommandRequirementsMockRecorder struct {
	mock *MockCommandRequirements
}

// NewMockCommandRequirements creates a new mock instance.
func NewMockCommandRequirements(ctrl *gomock.Controller) *MockCommandRequirements {
	mock := &MockCommandRequirements{ctrl: ctrl}
	mock.recorder = &MockCommandRequirementsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandRequirements) EXPECT() *MockCommandRequirementsMockRecorder {
	return m.recorder
}

// UsecaseQuery mocks base method.
func (m *MockCommandRequirements) UsecaseQuery() query.UsecaseQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsecaseQuery")
	ret0, _ := ret[0].(query.UsecaseQuery)
	return ret0
}

// UsecaseQuery indicates an expected call of UsecaseQuery.
func (mr *MockCommandRequirementsMockRecorder) UsecaseQuery() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsecaseQuery", reflect.TypeOf((*MockCommandRequirements)(nil).UsecaseQuery))
}
