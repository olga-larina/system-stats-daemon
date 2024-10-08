// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MetricCollector is an autogenerated mock type for the MetricCollector type
type MetricCollector struct {
	mock.Mock
}

type MetricCollector_Expecter struct {
	mock *mock.Mock
}

func (_m *MetricCollector) EXPECT() *MetricCollector_Expecter {
	return &MetricCollector_Expecter{mock: &_m.Mock}
}

// ExecuteCommand provides a mock function with given fields:
func (_m *MetricCollector) ExecuteCommand() ([]byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ExecuteCommand")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MetricCollector_ExecuteCommand_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteCommand'
type MetricCollector_ExecuteCommand_Call struct {
	*mock.Call
}

// ExecuteCommand is a helper method to define mock.On call
func (_e *MetricCollector_Expecter) ExecuteCommand() *MetricCollector_ExecuteCommand_Call {
	return &MetricCollector_ExecuteCommand_Call{Call: _e.mock.On("ExecuteCommand")}
}

func (_c *MetricCollector_ExecuteCommand_Call) Run(run func()) *MetricCollector_ExecuteCommand_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetricCollector_ExecuteCommand_Call) Return(_a0 []byte, _a1 error) *MetricCollector_ExecuteCommand_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MetricCollector_ExecuteCommand_Call) RunAndReturn(run func() ([]byte, error)) *MetricCollector_ExecuteCommand_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *MetricCollector) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MetricCollector_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MetricCollector_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MetricCollector_Expecter) Name() *MetricCollector_Name_Call {
	return &MetricCollector_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MetricCollector_Name_Call) Run(run func()) *MetricCollector_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetricCollector_Name_Call) Return(_a0 string) *MetricCollector_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetricCollector_Name_Call) RunAndReturn(run func() string) *MetricCollector_Name_Call {
	_c.Call.Return(run)
	return _c
}

// ParseCommandOutput provides a mock function with given fields: output
func (_m *MetricCollector) ParseCommandOutput(output string) (interface{}, error) {
	ret := _m.Called(output)

	if len(ret) == 0 {
		panic("no return value specified for ParseCommandOutput")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interface{}, error)); ok {
		return rf(output)
	}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(output)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(output)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MetricCollector_ParseCommandOutput_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseCommandOutput'
type MetricCollector_ParseCommandOutput_Call struct {
	*mock.Call
}

// ParseCommandOutput is a helper method to define mock.On call
//   - output string
func (_e *MetricCollector_Expecter) ParseCommandOutput(output interface{}) *MetricCollector_ParseCommandOutput_Call {
	return &MetricCollector_ParseCommandOutput_Call{Call: _e.mock.On("ParseCommandOutput", output)}
}

func (_c *MetricCollector_ParseCommandOutput_Call) Run(run func(output string)) *MetricCollector_ParseCommandOutput_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MetricCollector_ParseCommandOutput_Call) Return(_a0 interface{}, _a1 error) *MetricCollector_ParseCommandOutput_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MetricCollector_ParseCommandOutput_Call) RunAndReturn(run func(string) (interface{}, error)) *MetricCollector_ParseCommandOutput_Call {
	_c.Call.Return(run)
	return _c
}

// NewMetricCollector creates a new instance of MetricCollector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMetricCollector(t interface {
	mock.TestingT
	Cleanup(func())
}) *MetricCollector {
	mock := &MetricCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
