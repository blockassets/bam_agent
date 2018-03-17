package os

import "testing"

func TestReboot_Reboot(t *testing.T) {
	expectedCmd := "/sbin/reboot"
	expectedArg := "-f"

	reboot := &RebootData{
		run: func(cmd string, arg string) error {
			if cmd != expectedCmd {
				t.Fatalf("expected %s, got %s", expectedCmd, cmd)
			}
			if arg != expectedArg {
				t.Fatalf("expected %s, got %s", expectedArg, arg)
			}
			return nil
		},
	}
	reboot.Reboot()
}
