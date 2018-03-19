package os

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"go.uber.org/fx"
)

const (
	uptimePath = "/proc/uptime"
)

type UptimeResultFunc func() UptimeResult

type UptimeResult struct {
	time.Duration
	Error error
}

func NewUptimeResultFunc() UptimeResultFunc {
	return func() UptimeResult {
		return getUptimeResult()
	}
}

func getUptimeResult() UptimeResult {
	data, err := ioutil.ReadFile(uptimePath)
	if err != nil {
		return UptimeResult{time.Duration(0), err}
	}
	return parseUptime(string(data))
}

func parseUptime(data string) UptimeResult {
	parts := strings.Fields(data)
	if len(parts) < 2 {
		return UptimeResult{time.Duration(0), fmt.Errorf("unexpected content in %s: %s", uptimePath, data)}
	}

	uptimeInSeconds := strings.Split(parts[0], ".")
	if len(uptimeInSeconds) < 2 {
		return UptimeResult{time.Duration(0), fmt.Errorf("no period found %s", parts[0])}
	}

	flt, err := strconv.ParseInt(uptimeInSeconds[0], 10, 64)
	if err != nil {
		return UptimeResult{time.Duration(0), err}
	}

	return UptimeResult{time.Duration(flt) * time.Second, nil}
}

var UptimeModule = fx.Provide(NewUptimeResultFunc)
