// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/dappley/go-dappley/core"
import mock "github.com/stretchr/testify/mock"

// Consensus is an autogenerated mock type for the Consensus type
type Consensus struct {
	mock.Mock
}

// AddProducer provides a mock function with given fields: _a0
func (_m *Consensus) AddProducer(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FinishedMining provides a mock function with given fields:
func (_m *Consensus) FinishedMining() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetProducers provides a mock function with given fields:
func (_m *Consensus) GetProducers() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// SetKey provides a mock function with given fields: _a0
func (_m *Consensus) SetKey(_a0 string) {
	_m.Called(_a0)
}

// SetTargetBit provides a mock function with given fields: _a0
func (_m *Consensus) SetTargetBit(_a0 int) {
	_m.Called(_a0)
}

// Setup provides a mock function with given fields: _a0, _a1
func (_m *Consensus) Setup(_a0 core.NetService, _a1 string) {
	_m.Called(_a0, _a1)
}

// Start provides a mock function with given fields:
func (_m *Consensus) Start() {
	_m.Called()
}

// StartNewBlockMinting provides a mock function with given fields:
func (_m *Consensus) StartNewBlockMinting() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *Consensus) Stop() {
	_m.Called()
}

// Validate provides a mock function with given fields: block
func (_m *Consensus) Validate(block *core.Block) bool {
	ret := _m.Called(block)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*core.Block) bool); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// VerifyBlock provides a mock function with given fields: block
func (_m *Consensus) VerifyBlock(block *core.Block) bool {
	ret := _m.Called(block)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*core.Block) bool); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
