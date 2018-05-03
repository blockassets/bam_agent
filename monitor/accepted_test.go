package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
)

func TestNewAcceptedMonitor(t *testing.T) {
	config := AcceptedConfig{Enabled: true, Period: time.Duration(5)}

	client := miner.NewMockMinerClient(miner.AcceptedIncrement)
	reboot := os.NewMockReboot()
	result := NewAcceptedMonitor(config, &client, &reboot)
	monitor := result.Monitor

	if monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected period 5, got %s", monitor.GetPeriod())
	}

	if !monitor.IsEnabled() {
		t.Fatalf("expected enabled, got %v", monitor.IsEnabled())
	}

	// Run
	handlerFunc := result.Monitor.NewTickerFunc()
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	client = miner.NewMockMinerClient(miner.AcceptedError)
	result = NewAcceptedMonitor(config, &client, &reboot)
	monitor = result.Monitor
	handlerFunc = result.Monitor.NewTickerFunc()
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	client = miner.NewMockMinerClient(miner.AcceptedSame)
	result = NewAcceptedMonitor(config, &client, &reboot)
	monitor = result.Monitor
	handlerFunc = result.Monitor.NewTickerFunc()
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	if !reboot.CalledReboot {
		t.Fatalf("expected reboot, go %v", reboot.CalledReboot)
	}
	if reboot.Counter != 2 {
		t.Fatalf("expected counter 2, got %d", reboot.Counter)
	}

	reboot = os.NewMockReboot()
	client = miner.NewMockMinerClient(miner.AcceptedZero)
	result = NewAcceptedMonitor(config, &client, &reboot)
	monitor = result.Monitor
	handlerFunc = result.Monitor.NewTickerFunc()
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	client = miner.NewMockMinerClient(miner.AcceptedIncrement)
	result = NewAcceptedMonitor(config, &client, &reboot)
	monitor = result.Monitor
	handlerFunc = result.Monitor.NewTickerFunc()
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	handlerFunc(context.TODO())
	client.Test = miner.AcceptedError
	handlerFunc(context.TODO())
	if !reboot.CalledReboot {
		t.Fatalf("expected reboot, got %v", reboot.CalledReboot)
	}
}
