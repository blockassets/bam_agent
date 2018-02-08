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
	getLoad() (loads []float64, err error)
type statRetrieve interface {
}

type LinuxStatRetriever struct {
}

// DefaultMountPoint is the common mount point of the proc filesystem.
const DefaultMountPoint = "/proc"

func procFilePath(name string) string {
	return path.Join(DefaultMountPoint, name)
}

// Read loadavg from /proc.
func (LinuxStatRetriever) getLoad() (loads []float64, err error) {
	data, err := ioutil.ReadFile(procFilePath("loadavg"))
	if err != nil {
		return nil, err
	}
	loads, err = parseLoad(string(data))
	if err != nil {
		return nil, err
	}
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
