package monitor

import "time"

type PeriodicAction interface {
	ReloadConfig(blob []byte)
	IsEnabled() bool
	InitialPeriod() time.Duration
	Period() time.Duration
	Action()
}

func periodicMonitor(msg_ch chan interface{}, maxVelocity time.Duration, pa PeriodicAction) {

	paused := false
	// We dont use a select here as it seems awkward to manage turning on and off the timer
	// so use the channel to receieve both timer messages and also broadcast messages

	for m := range msg_ch {
		msgR := m.(msg)
		switch msgR.msgType {
		case msgReloadConfig:
			{
				pa.ReloadConfig(msgR.msgBody)
				if pa.IsEnabled() {
					if pa.InitialPeriod() > maxVelocity {
						sendMsgTimerAfter(msg_ch, pa.InitialPeriod())
					}
				}
			}
		case msgTimer:
			{
				if pa.IsEnabled() {
					if !paused {
						pa.Action()
					}
					if pa.Period() > maxVelocity {
						sendMsgTimerAfter(msg_ch, pa.Period())
					}

				}
			}
		case msgPause:
			{
				paused = true
			}
		case msgUnpause:
			{
				paused = false
			}
		}

	}
}
