package os

import (
	"io"
)

const (
	MemBelow = iota
	MemExactly
	MemAbove
)

// Type insurance
var _ MemInfo = &MockMemInfo{}

type MockMemInfo struct {
	dataSet   int
	path      string
	getReader func(path string) (io.Reader, error)
}

func (sr *MockMemInfo) Get() (MemInfoData, error) {
	memInfo := MemInfoData{}
	switch sr.dataSet {
	case MemBelow:
		memInfo[MemAvailable] = 100 * 1000 * 1000 // bytes
	case MemExactly:
		memInfo[MemAvailable] = 125 * 1000 * 1000 // bytes
	case MemAbove:
		memInfo[MemAvailable] = 400 * 1000 * 1000 // bytes
	}

	return memInfo, nil
}

func NewMockMemInfo(dataSet int) MockMemInfo {
	return MockMemInfo{
		dataSet: dataSet,
	}
}
