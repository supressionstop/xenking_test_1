// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// FetchLiner is an autogenerated mock type for the FetchLiner type
type FetchLiner struct {
	mock.Mock
}

type FetchLiner_Expecter struct {
	mock *mock.Mock
}

func (_m *FetchLiner) EXPECT() *FetchLiner_Expecter {
	return &FetchLiner_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, sport
func (_m *FetchLiner) Execute(ctx context.Context, sport string) error {
	ret := _m.Called(ctx, sport)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, sport)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchLiner_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type FetchLiner_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - sport string
func (_e *FetchLiner_Expecter) Execute(ctx interface{}, sport interface{}) *FetchLiner_Execute_Call {
	return &FetchLiner_Execute_Call{Call: _e.mock.On("Execute", ctx, sport)}
}

func (_c *FetchLiner_Execute_Call) Run(run func(ctx context.Context, sport string)) *FetchLiner_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *FetchLiner_Execute_Call) Return(_a0 error) *FetchLiner_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FetchLiner_Execute_Call) RunAndReturn(run func(context.Context, string) error) *FetchLiner_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewFetchLiner creates a new instance of FetchLiner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFetchLiner(t interface {
	mock.TestingT
	Cleanup(func())
}) *FetchLiner {
	mock := &FetchLiner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
