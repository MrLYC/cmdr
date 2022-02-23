// Code generated by MockGen. DO NOT EDIT.
// Source: command.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	core "github.com/mrlyc/cmdr/core"
)

// MockCommand is a mock of Command interface.
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand.
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance.
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// GetActivated mocks base method.
func (m *MockCommand) GetActivated() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivated")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetActivated indicates an expected call of GetActivated.
func (mr *MockCommandMockRecorder) GetActivated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivated", reflect.TypeOf((*MockCommand)(nil).GetActivated))
}

// GetLocation mocks base method.
func (m *MockCommand) GetLocation() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocation")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLocation indicates an expected call of GetLocation.
func (mr *MockCommandMockRecorder) GetLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocation", reflect.TypeOf((*MockCommand)(nil).GetLocation))
}

// GetName mocks base method.
func (m *MockCommand) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockCommandMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockCommand)(nil).GetName))
}

// GetVersion mocks base method.
func (m *MockCommand) GetVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockCommandMockRecorder) GetVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockCommand)(nil).GetVersion))
}

// MockCommandQuery is a mock of CommandQuery interface.
type MockCommandQuery struct {
	ctrl     *gomock.Controller
	recorder *MockCommandQueryMockRecorder
}

// MockCommandQueryMockRecorder is the mock recorder for MockCommandQuery.
type MockCommandQueryMockRecorder struct {
	mock *MockCommandQuery
}

// NewMockCommandQuery creates a new mock instance.
func NewMockCommandQuery(ctrl *gomock.Controller) *MockCommandQuery {
	mock := &MockCommandQuery{ctrl: ctrl}
	mock.recorder = &MockCommandQueryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandQuery) EXPECT() *MockCommandQueryMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockCommandQuery) All() ([]core.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All")
	ret0, _ := ret[0].([]core.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockCommandQueryMockRecorder) All() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockCommandQuery)(nil).All))
}

// Count mocks base method.
func (m *MockCommandQuery) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockCommandQueryMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockCommandQuery)(nil).Count))
}

// One mocks base method.
func (m *MockCommandQuery) One() (core.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "One")
	ret0, _ := ret[0].(core.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// One indicates an expected call of One.
func (mr *MockCommandQueryMockRecorder) One() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "One", reflect.TypeOf((*MockCommandQuery)(nil).One))
}

// WithActivated mocks base method.
func (m *MockCommandQuery) WithActivated(activated bool) core.CommandQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithActivated", activated)
	ret0, _ := ret[0].(core.CommandQuery)
	return ret0
}

// WithActivated indicates an expected call of WithActivated.
func (mr *MockCommandQueryMockRecorder) WithActivated(activated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithActivated", reflect.TypeOf((*MockCommandQuery)(nil).WithActivated), activated)
}

// WithLocation mocks base method.
func (m *MockCommandQuery) WithLocation(location string) core.CommandQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithLocation", location)
	ret0, _ := ret[0].(core.CommandQuery)
	return ret0
}

// WithLocation indicates an expected call of WithLocation.
func (mr *MockCommandQueryMockRecorder) WithLocation(location interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLocation", reflect.TypeOf((*MockCommandQuery)(nil).WithLocation), location)
}

// WithName mocks base method.
func (m *MockCommandQuery) WithName(name string) core.CommandQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithName", name)
	ret0, _ := ret[0].(core.CommandQuery)
	return ret0
}

// WithName indicates an expected call of WithName.
func (mr *MockCommandQueryMockRecorder) WithName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithName", reflect.TypeOf((*MockCommandQuery)(nil).WithName), name)
}

// WithVersion mocks base method.
func (m *MockCommandQuery) WithVersion(version string) core.CommandQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithVersion", version)
	ret0, _ := ret[0].(core.CommandQuery)
	return ret0
}

// WithVersion indicates an expected call of WithVersion.
func (mr *MockCommandQueryMockRecorder) WithVersion(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithVersion", reflect.TypeOf((*MockCommandQuery)(nil).WithVersion), version)
}

// MockCommandManager is a mock of CommandManager interface.
type MockCommandManager struct {
	ctrl     *gomock.Controller
	recorder *MockCommandManagerMockRecorder
}

// MockCommandManagerMockRecorder is the mock recorder for MockCommandManager.
type MockCommandManagerMockRecorder struct {
	mock *MockCommandManager
}

// NewMockCommandManager creates a new mock instance.
func NewMockCommandManager(ctrl *gomock.Controller) *MockCommandManager {
	mock := &MockCommandManager{ctrl: ctrl}
	mock.recorder = &MockCommandManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandManager) EXPECT() *MockCommandManagerMockRecorder {
	return m.recorder
}

// Activate mocks base method.
func (m *MockCommandManager) Activate(name, version string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Activate", name, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Activate indicates an expected call of Activate.
func (mr *MockCommandManagerMockRecorder) Activate(name, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Activate", reflect.TypeOf((*MockCommandManager)(nil).Activate), name, version)
}

// Close mocks base method.
func (m *MockCommandManager) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockCommandManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCommandManager)(nil).Close))
}

// Deactivate mocks base method.
func (m *MockCommandManager) Deactivate(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deactivate", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deactivate indicates an expected call of Deactivate.
func (mr *MockCommandManagerMockRecorder) Deactivate(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deactivate", reflect.TypeOf((*MockCommandManager)(nil).Deactivate), name)
}

// Define mocks base method.
func (m *MockCommandManager) Define(name, version, location string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Define", name, version, location)
	ret0, _ := ret[0].(error)
	return ret0
}

// Define indicates an expected call of Define.
func (mr *MockCommandManagerMockRecorder) Define(name, version, location interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Define", reflect.TypeOf((*MockCommandManager)(nil).Define), name, version, location)
}

// Provider mocks base method.
func (m *MockCommandManager) Provider() core.CommandProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Provider")
	ret0, _ := ret[0].(core.CommandProvider)
	return ret0
}

// Provider indicates an expected call of Provider.
func (mr *MockCommandManagerMockRecorder) Provider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Provider", reflect.TypeOf((*MockCommandManager)(nil).Provider))
}

// Query mocks base method.
func (m *MockCommandManager) Query() (core.CommandQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query")
	ret0, _ := ret[0].(core.CommandQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockCommandManagerMockRecorder) Query() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockCommandManager)(nil).Query))
}

// Undefine mocks base method.
func (m *MockCommandManager) Undefine(name, version string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Undefine", name, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Undefine indicates an expected call of Undefine.
func (mr *MockCommandManagerMockRecorder) Undefine(name, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Undefine", reflect.TypeOf((*MockCommandManager)(nil).Undefine), name, version)
}