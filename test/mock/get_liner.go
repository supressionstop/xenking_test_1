// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/supressionstop/xenking_test_1/internal/usecase/entity"
)

// GetLiner is an autogenerated mock type for the GetLiner type
type GetLiner struct {
	mock.Mock
}

type GetLiner_Expecter struct {
	mock *mock.Mock
}

func (_m *GetLiner) EXPECT() *GetLiner_Expecter {
	return &GetLiner_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, sport
func (_m *GetLiner) Execute(ctx context.Context, sport string) (entity.Line, error) {
	ret := _m.Called(ctx, sport)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 entity.Line
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (entity.Line, error)); ok {
		return rf(ctx, sport)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.Line); ok {
		r0 = rf(ctx, sport)
	} else {
		r0 = ret.Get(0).(entity.Line)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, sport)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLiner_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type GetLiner_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - sport string
func (_e *GetLiner_Expecter) Execute(ctx interface{}, sport interface{}) *GetLiner_Execute_Call {
	return &GetLiner_Execute_Call{Call: _e.mock.On("Execute", ctx, sport)}
}

func (_c *GetLiner_Execute_Call) Run(run func(ctx context.Context, sport string)) *GetLiner_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *GetLiner_Execute_Call) Return(_a0 entity.Line, _a1 error) *GetLiner_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GetLiner_Execute_Call) RunAndReturn(run func(context.Context, string) (entity.Line, error)) *GetLiner_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewGetLiner creates a new instance of GetLiner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGetLiner(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetLiner {
	mock := &GetLiner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
