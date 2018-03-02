package service

// this code cut and pasted from the prometheus repo

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// facilitate testing of code that uses this package
type StatRetriever interface {
	GetLoad() (loads LoadAvgs, err error)
}

type LinuxStatRetriever struct {
}

func NewStatRetriever() StatRetriever {
	return &LinuxStatRetriever{}
}

type LoadAvgs struct {
	OneMinAvg     float64
	FiveMinAvg    float64
	FifteenMinAvg float64
}

// DefaultMountPoint is the common mount point of the proc filesystem.
const loadAvgProc = "/proc/loadavg"

// Read loadavg from /proc.
func (*LinuxStatRetriever) GetLoad() (LoadAvgs, error) {
	data, err := ioutil.ReadFile(loadAvgProc)
	if err != nil {
		return LoadAvgs{}, err
	}
	return ParseLoad(string(data))
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func ParseLoad(data string) (LoadAvgs, error) {
	loadsAsArray := make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return LoadAvgs{}, fmt.Errorf("unexpected content %s", data)
	}
	var err error
	for i, load := range parts[0:3] {
		loadsAsArray[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return LoadAvgs{}, fmt.Errorf("could not parse load '%v': %v", load, err)
		}
	}
	loads := LoadAvgs{
		OneMinAvg:     loadsAsArray[0],
		FiveMinAvg:    loadsAsArray[1],
		FifteenMinAvg: loadsAsArray[2],
	}
	return loads, nil
}
