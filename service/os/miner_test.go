package os

import "testing"

/*
	The best we can really do is unit level testing in here.
	Integration tests would pick up that os.exec is called.
*/
func TestSystem_StartMiner(t *testing.T) {
	expected := "systemctl start cgminer"
	miner := &MinerData{
		run: func(cmd string) error {
			if cmd != expected {
				t.Fatalf("expected %s, got %s", expected, cmd)
			}
			return nil
		},
	}
	miner.Start()
}

func TestSystem_StopMiner(t *testing.T) {
	expected := "systemctl stop cgminer"
	miner := &MinerData{
		run: func(cmd string) error {
			if cmd != expected {
				t.Fatalf("expected %s, got %s", expected, cmd)
			}
			return nil
		},
	}
	miner.Stop()
}
