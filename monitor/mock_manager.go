package monitor

// Type insurance
var _ Manager = &MockManager{}

type MockManager struct {
	CalledStart bool
	CalledStop  bool
}

func (mm *MockManager) Start() {
	mm.CalledStart = true
}

func (mm *MockManager) Stop() {
	mm.CalledStop = true
}

func NewMockManager() MockManager {
	return MockManager{}
}
