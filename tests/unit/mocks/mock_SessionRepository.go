// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockSessionRepository is an autogenerated mock type for the SessionRepository type
type MockSessionRepository struct {
	mock.Mock
}

type MockSessionRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSessionRepository) EXPECT() *MockSessionRepository_Expecter {
	return &MockSessionRepository_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: token
func (_m *MockSessionRepository) Add(token string) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSessionRepository_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockSessionRepository_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - token string
func (_e *MockSessionRepository_Expecter) Add(token interface{}) *MockSessionRepository_Add_Call {
	return &MockSessionRepository_Add_Call{Call: _e.mock.On("Add", token)}
}

func (_c *MockSessionRepository_Add_Call) Run(run func(token string)) *MockSessionRepository_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockSessionRepository_Add_Call) Return(_a0 error) *MockSessionRepository_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSessionRepository_Add_Call) RunAndReturn(run func(string) error) *MockSessionRepository_Add_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSessionRepository creates a new instance of MockSessionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSessionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSessionRepository {
	mock := &MockSessionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}