// Code generated by mockery v2.19.0. DO NOT EDIT.

package mocks

import (
	context "context"

	client "sigs.k8s.io/controller-runtime/pkg/client"

	mock "github.com/stretchr/testify/mock"

	model "github.com/NetApp-Polaris/astra-connector-operator/app/deployer/model"

	v1 "github.com/NetApp-Polaris/astra-connector-operator/details/operator-sdk/api/v1"
)

// getK8sResources is an autogenerated mock type for the getK8sResources type
type getK8sResources struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *getK8sResources) Execute(_a0 model.Deployer, _a1 *v1.AstraConnector, _a2 context.Context) ([]client.Object, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []client.Object
	if rf, ok := ret.Get(0).(func(model.Deployer, *v1.AstraConnector, context.Context) []client.Object); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]client.Object)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.Deployer, *v1.AstraConnector, context.Context) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewGetK8sResources interface {
	mock.TestingT
	Cleanup(func())
}

// newGetK8sResources creates a new instance of getK8sResources. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newGetK8sResources(t mockConstructorTestingTnewGetK8sResources) *getK8sResources {
	mock := &getK8sResources{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
