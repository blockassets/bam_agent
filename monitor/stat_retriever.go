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
func (LinuxStatRetriever) getLoad() (LoadAvgs, error) {
	data, err := ioutil.ReadFile(procFilePath("loadavg"))
	if err != nil {
		return LoadAvgs{}, err
	}
	return parseLoad(string(data))
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func parseLoad(data string) (LoadAvgs, error) {
	loadsAsArray := make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return LoadAvgs{}, fmt.Errorf("unexpected content in %v", procFilePath("loadavg"))
	}
	var err error
	for i, load := range parts[0:3] {
		loadsAsArray[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return LoadAvgs{}, fmt.Errorf("could not parse load '%v': %v", load, err)
		}
	}
	loads := LoadAvgs{
		oneMinAvg:     loadsAsArray[0],
		fiveMinAvg:    loadsAsArray[1],
		fifteenMinAvg: loadsAsArray[2],
	}
	return loads, nil
}
