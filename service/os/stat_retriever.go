package os

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"go.uber.org/fx"
)

// Common on most systems
const defaultLoadPath = "/proc/loadavg"

type StatRetriever interface {
	GetLoadData() (data *LoadData, err error)
}

type LinuxStatRetriever struct {
	loadPath    string
	getProcData func(loadPath string) ([]byte, error)
}

type LoadData struct {
	OneMinAvg     float64
	FiveMinAvg    float64
	FifteenMinAvg float64
}

func (lsr *LinuxStatRetriever) GetLoadData() (*LoadData, error) {
	data, err := lsr.getProcData(lsr.loadPath)
	if err != nil {
		return nil, err
	}
	return ParseLoadData(string(data))
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func ParseLoadData(data string) (*LoadData, error) {
	loadsAsArray := make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content %s", data)
	}
	var err error
	for i, load := range parts[0:3] {
		loadsAsArray[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse load '%v': %v", load, err)
		}
	}
	loads := &LoadData{
		OneMinAvg:     loadsAsArray[0],
		FiveMinAvg:    loadsAsArray[1],
		FifteenMinAvg: loadsAsArray[2],
	}
	return loads, nil
}

func NewStatRetriever() StatRetriever {
	return &LinuxStatRetriever{
		loadPath: defaultLoadPath,
		getProcData: func(loadPath string) ([]byte, error) {
			return ioutil.ReadFile(loadPath)
		},
	}
}

var StatRetrieverModule = fx.Provide(NewStatRetriever)
