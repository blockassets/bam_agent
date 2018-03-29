package os

import "testing"

func TestReboot_Reboot(t *testing.T) {
	expectedRebootCmd := "/sbin/reboot"
	expectedRebootArg := "-f"

	expectedSyncCmd := "/bin/sync"
	syncCounter := 0

	reboot := &RebootData{
		run: func(cmd string, arg string) error {
			if arg != "" {
				if cmd != expectedRebootCmd {
					t.Fatalf("expected %s, got %s", expectedRebootCmd, cmd)
				}
				if arg != expectedRebootArg {
					t.Fatalf("expected %s, got %s", expectedRebootArg, arg)
				}
			} else {
				if cmd != expectedSyncCmd {
					t.Fatalf("expected %s, got %s", expectedSyncCmd, cmd)
				}
				syncCounter++
			}

			return nil
		},
	}
	reboot.Reboot()

	if syncCounter != 2 {
		t.Fatalf("expected sync counter to be 2, got %v", syncCounter)
	}
}
