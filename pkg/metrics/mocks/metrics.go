// Code generated by MockGen. DO NOT EDIT.
// Source: metrics.go

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetrics is a mock of Metrics interface.
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics.
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance.
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// IncBytesUploaded mocks base method.
func (m *MockMetrics) IncBytesUploaded(bytesUploaded int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncBytesUploaded", bytesUploaded)
}

// IncBytesUploaded indicates an expected call of IncBytesUploaded.
func (mr *MockMetricsMockRecorder) IncBytesUploaded(bytesUploaded interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncBytesUploaded", reflect.TypeOf((*MockMetrics)(nil).IncBytesUploaded), bytesUploaded)
}

// IncHits mocks base method.
func (m *MockMetrics) IncHits(status int, method, path string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncHits", status, method, path)
}

// IncHits indicates an expected call of IncHits.
func (mr *MockMetricsMockRecorder) IncHits(status, method, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncHits", reflect.TypeOf((*MockMetrics)(nil).IncHits), status, method, path)
}

// ObserveResponseTime mocks base method.
func (m *MockMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ObserveResponseTime", status, method, path, observeTime)
}

// ObserveResponseTime indicates an expected call of ObserveResponseTime.
func (mr *MockMetricsMockRecorder) ObserveResponseTime(status, method, path, observeTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObserveResponseTime", reflect.TypeOf((*MockMetrics)(nil).ObserveResponseTime), status, method, path, observeTime)
}
