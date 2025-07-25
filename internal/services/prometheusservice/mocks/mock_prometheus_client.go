// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/motain/of-catalog/internal/services/prometheusservice (interfaces: PrometheusClientInterface)

// Package prometheusservice is a generated GoMock package.
package prometheusservice

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	model "github.com/prometheus/common/model"
)

// MockPrometheusClientInterface is a mock of PrometheusClientInterface interface.
type MockPrometheusClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPrometheusClientInterfaceMockRecorder
}

// MockPrometheusClientInterfaceMockRecorder is the mock recorder for MockPrometheusClientInterface.
type MockPrometheusClientInterfaceMockRecorder struct {
	mock *MockPrometheusClientInterface
}

// NewMockPrometheusClientInterface creates a new mock instance.
func NewMockPrometheusClientInterface(ctrl *gomock.Controller) *MockPrometheusClientInterface {
	mock := &MockPrometheusClientInterface{ctrl: ctrl}
	mock.recorder = &MockPrometheusClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrometheusClientInterface) EXPECT() *MockPrometheusClientInterfaceMockRecorder {
	return m.recorder
}

// Query mocks base method.
func (m *MockPrometheusClientInterface) Query(arg0 string, arg1 time.Time) (model.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", arg0, arg1)
	ret0, _ := ret[0].(model.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockPrometheusClientInterfaceMockRecorder) Query(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPrometheusClientInterface)(nil).Query), arg0, arg1)
}

// QueryRange mocks base method.
func (m *MockPrometheusClientInterface) QueryRange(arg0 string, arg1 v1.Range) (model.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryRange", arg0, arg1)
	ret0, _ := ret[0].(model.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryRange indicates an expected call of QueryRange.
func (mr *MockPrometheusClientInterfaceMockRecorder) QueryRange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRange", reflect.TypeOf((*MockPrometheusClientInterface)(nil).QueryRange), arg0, arg1)
}
