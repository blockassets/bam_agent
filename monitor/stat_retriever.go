package monitor

// this code cut and pasted from the promethious repo

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

// facilitate testing of code that uses this package
type statRetriever interface {
	getLoad() (loads LoadAvgs, err error)
}

type LinuxStatRetriever struct {
}

type LoadAvgs struct {
	oneMinAvg     float64
	fiveMinAvg    float64
	fifteenMinAvg float64
}

// DefaultMountPoint is the common mount point of the proc filesystem.
const DefaultMountPoint = "/proc"

func procFilePath(name string) string {
	return path.Join(DefaultMountPoint, name)
}

// Read loadavg from /proc.
func (LinuxStatRetriever) getLoad() (loads LoadAvgs, err error) {
	data, err := ioutil.ReadFile(procFilePath("loadavg"))
	if err != nil {
		return loads, err
	}

	loadsAsArray, err := parseLoad(string(data))
	if err != nil {
		return loads, err
	}
	loads.oneMinAvg = loadsAsArray[0]
	loads.fiveMinAvg = loadsAsArray[1]
	loads.fifteenMinAvg = loadsAsArray[2]

	return loads, nil
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func parseLoad(data string) (loads []float64, err error) {
	loads = make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content in %v", procFilePath("loadavg"))
	}
	for i, load := range parts[0:3] {
		loads[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse load '%v': %v", load, err)
		}
	}
	return loads, nil
}
