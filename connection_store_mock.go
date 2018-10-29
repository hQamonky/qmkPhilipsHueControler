package main

import (
	"github.com/stretchr/testify/mock"
)

// MockStore contains additonal methods for inspection
type MockStore struct {
	mock.Mock
}

// CreateBridge is
func (m *MockStore) CreateBridge(bridge *Bridge) error {
	/*
		When this method is called, `m.Called` records the call, and also
		returns the result that we pass to it (which you will see in the
		handler tests)
	*/
	rets := m.Called(bridge)
	return rets.Error(0)
}

// GetBridges is
func (m *MockStore) GetBridges() ([]*Bridge, error) {
	rets := m.Called()
	/*
		Since `rets.Get()` is a generic method, that returns whatever we pass to it,
		we need to typecast it to the type we expect, which in this case is []*Bridge
	*/
	return rets.Get(0).([]*Bridge), rets.Error(1)
}

// InitMockStore is
func InitMockStore() *MockStore {
	/*
		Like the InitStore function we defined earlier, this function
		also initializes the store variable, but this time, it assigns
		a new MockStore instance to it, instead of an actual store
	*/
	s := new(MockStore)
	store = s
	return s
}
