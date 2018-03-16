package os

const (
	LevelNotEnough = iota
	LevelBelowFive
	LevelExactlyFive
	LevelAboveFive
	LevelMalformed
)

// Type insurance
var _ StatRetriever = &MockStatRetriever{}

type MockStatRetriever struct {
	dataSet     int
	loadPath    string
	getProcData func(loadPath string) ([]byte, error)
}

func (sr *MockStatRetriever) GetLoadData() (*LoadData, error) {
	pd, err := sr.getProcData("")
	if err != nil {
		return nil, err
	}
	return ParseLoadData(string(pd))
}

func NewMockStatRetriever(dataSet int) MockStatRetriever {
	var data string
	switch dataSet {
	case LevelNotEnough:
		data = "0.0 0.0"
	case LevelBelowFive:
		data = "0.0 4.999 0.0 1234 1234"
	case LevelExactlyFive:
		data = "0.0 5.0 0.0 1234 1234"
	case LevelAboveFive:
		data = "0.0 5.1 0.0 1234 1234"
	case LevelMalformed:
		data = "a b c d emnf,masfd"
	}

	return MockStatRetriever{
		getProcData: func(loadPath string) ([]byte, error) {
			return []byte(data), nil
		},
	}
}
