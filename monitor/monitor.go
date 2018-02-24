package monitor

import (
	"log"
	"math/rand"
	"time"

	"github.com/blockassets/bam_agent/util"
)

const (
	msgReloadConfig = 1
	msgTimer        = 2
	msgPause        = 3
	msgUnpause      = 4
)

type msg struct {
	msgType int
	msgBody []byte
}

//common attributes for the PeriodicAction module
type MonitorConfig struct {
	Enabled                bool `json:"enabled"`
	PeriodSecs             int  `json:"period_secs"`
	InitialPeriodRangeSecs int  `json:"initial_range_secs"`
}

//default handlers for the PeriodicAction interface
func (cfg *MonitorConfig) IsEnabled() bool {
	return cfg.Enabled
}

func (cfg *MonitorConfig) InitialPeriod() time.Duration {
	// Use a random number to random set the inital period. So if a 1000 miners are reset, they come back on line in a random
	// distribution so that we dont get seen as a denial of service attack on the pool
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return time.Duration(cfg.PeriodSecs)*time.Second + time.Duration(r1.Intn(cfg.InitialPeriodRangeSecs))*time.Second
}

func (cfg *MonitorConfig) Period() time.Duration {
	return time.Duration(cfg.PeriodSecs) * time.Second
}

// package scope to hold the collection of channels used to communicate with the monitors
var monitorControl util.Broadcaster

func StartMonitors(c *util.ConfigFile) {

	sr := LinuxStatRetriever{}
	monitorControl = util.NewBroadcaster(0)

	log.Println("Monitors being started")
	ch := make(chan interface{})
	monitorControl.Register(ch)
	go monitorLoad(ch, sr) // check for 5min average CPU every so often

	ch = make(chan interface{})
	monitorControl.Register(ch)
	go monitorPeriodicMinerQuit(ch)

	ch = make(chan interface{})
	monitorControl.Register(ch)
	go monitorPeriodicReboot(ch)

	ReLoadMonitorConfiguration(c)
}

func ReLoadMonitorConfiguration(c *util.ConfigFile) {
	buf := c.GetConfigBuf()
	m := msg{msgReloadConfig, buf}
	monitorControl.Submit(m)
}

func PauseMonitors() {
	monitorControl.Submit(msg{msgPause, nil})

}

func UnPauseMonitors() {
	monitorControl.Submit(msg{msgUnpause, nil})
}

func sendMsgTimerAfter(ch chan interface{}, timeToWait time.Duration) {
	time.AfterFunc(timeToWait, func() {
		ch <- msg{msgTimer, nil}
		return
	})
}
