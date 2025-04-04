// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	accountmodel "raiseexception.dev/odin/src/accounting/domain/account"

	mock "github.com/stretchr/testify/mock"
)

// MockAccountRepository is an autogenerated mock type for the AccountRepository type
type MockAccountRepository struct {
	mock.Mock
}

type MockAccountRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccountRepository) EXPECT() *MockAccountRepository_Expecter {
	return &MockAccountRepository_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, account
func (_m *MockAccountRepository) Add(ctx context.Context, account *accountmodel.Account) error {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *accountmodel.Account) error); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccountRepository_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockAccountRepository_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - account *accountmodel.Account
func (_e *MockAccountRepository_Expecter) Add(ctx interface{}, account interface{}) *MockAccountRepository_Add_Call {
	return &MockAccountRepository_Add_Call{Call: _e.mock.On("Add", ctx, account)}
}

func (_c *MockAccountRepository_Add_Call) Run(run func(ctx context.Context, account *accountmodel.Account)) *MockAccountRepository_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*accountmodel.Account))
	})
	return _c
}

func (_c *MockAccountRepository_Add_Call) Return(_a0 error) *MockAccountRepository_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccountRepository_Add_Call) RunAndReturn(run func(context.Context, *accountmodel.Account) error) *MockAccountRepository_Add_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *MockAccountRepository) GetAll(ctx context.Context) ([]*accountmodel.Account, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*accountmodel.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*accountmodel.Account, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*accountmodel.Account); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*accountmodel.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockAccountRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAccountRepository_Expecter) GetAll(ctx interface{}) *MockAccountRepository_GetAll_Call {
	return &MockAccountRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *MockAccountRepository_GetAll_Call) Run(run func(ctx context.Context)) *MockAccountRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAccountRepository_GetAll_Call) Return(_a0 []*accountmodel.Account, _a1 error) *MockAccountRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountRepository_GetAll_Call) RunAndReturn(run func(context.Context) ([]*accountmodel.Account, error)) *MockAccountRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccountRepository creates a new instance of MockAccountRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccountRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccountRepository {
	mock := &MockAccountRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
