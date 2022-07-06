// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	storm "github.com/asdine/storm/v3"
	mock "github.com/stretchr/testify/mock"
)

// SaveStopper is an autogenerated mock type for the SaveStopper type
type SaveStopper struct {
	mock.Mock
}

// Begin provides a mock function with given fields: writable
func (_m *SaveStopper) Begin(writable bool) (storm.Node, error) {
	ret := _m.Called(writable)

	var r0 storm.Node
	if rf, ok := ret.Get(0).(func(bool) storm.Node); ok {
		r0 = rf(writable)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storm.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(writable)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: data
func (_m *SaveStopper) Save(data interface{}) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *SaveStopper) Stop() {
	_m.Called()
}

type mockConstructorTestingTNewSaveStopper interface {
	mock.TestingT
	Cleanup(func())
}

// NewSaveStopper creates a new instance of SaveStopper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSaveStopper(t mockConstructorTestingTNewSaveStopper) *SaveStopper {
	mock := &SaveStopper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}