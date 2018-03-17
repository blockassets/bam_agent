package os

import (
	"net"
)

// Type insurance
var _ NetInfo = &MockNetInfo{}

type MockNetInfo struct {
	data []net.Interface
}

func (mni *MockNetInfo) GetNetInterfaces() ([]net.Interface, error) {
	return mni.data, nil
}

func (mni *MockNetInfo) GetMacAddress() *string {
	str := mni.data[0].HardwareAddr.String()
	return &str
}

func NewMockNetInfo() MockNetInfo {
	return MockNetInfo{
		data: NewNetInterfaces(),
	}
}
