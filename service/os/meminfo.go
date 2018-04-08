package os

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/fx"
)

const defaultPath = "/proc/meminfo"

const MemAvailable = "MemAvailable_bytes"

type MemInfoData = map[string]float64

type MemInfo interface {
	Get() (data MemInfoData, err error)
}

type LinuxMemInfo struct {
	path    string
	getData func(path string) ([]byte, error)
}

func (m *LinuxMemInfo) Get() (MemInfoData, error) {
	data, err := m.getData(m.path)

	if err != nil {
		return nil, err
	}

	return parseMemInfo(bytes.NewReader(data))
}

func parseMemInfo(r io.Reader) (MemInfoData, error) {
	var (
		memInfo = map[string]float64{}
		scanner = bufio.NewScanner(r)
		re      = regexp.MustCompile(`\((.*)\)`)
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %s", err)
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// Active(anon) -> Active_anon
		key = re.ReplaceAllString(key, "_${1}")
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
			key = key + "_bytes"
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		memInfo[key] = fv
	}

	return memInfo, scanner.Err()
}

func NewMemInfo() MemInfo {
	return &LinuxMemInfo{
		path: defaultPath,
		getData: func(path string) ([]byte, error) {
			return ioutil.ReadFile(defaultPath)
		},
	}
}

var MemInfoModule = fx.Provide(NewMemInfo)
