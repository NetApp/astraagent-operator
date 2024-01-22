// Code generated by mockery v2.19.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ClusterTypeCheckerInterface is an autogenerated mock type for the ClusterTypeCheckerInterface type
type ClusterTypeCheckerInterface struct {
	mock.Mock
}

// DetermineClusterType provides a mock function with given fields:
func (_m *ClusterTypeCheckerInterface) DetermineClusterType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewClusterTypeCheckerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewClusterTypeCheckerInterface creates a new instance of ClusterTypeCheckerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClusterTypeCheckerInterface(t mockConstructorTestingTNewClusterTypeCheckerInterface) *ClusterTypeCheckerInterface {
	mock := &ClusterTypeCheckerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
