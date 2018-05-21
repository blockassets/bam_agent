package os

import "testing"

func TestNtpdate_Ntpdate(t *testing.T) {
	expectedNtpdateCmd := "/usr/bin/ntpdate"
	expectedNtpdateArgs := []string{"-u", "time.foo.com"}

	ntpdate := &NtpdateData{
		run: func(cmd string, arg ...string) error {
			if cmd != expectedNtpdateCmd {
				t.Fatalf("expected %s, got %s", expectedNtpdateCmd, cmd)
			}

			if arg[0] != expectedNtpdateArgs[0] {
				t.Fatalf("expected %s, got %s", expectedNtpdateArgs[0], arg[0])
			}

			if arg[1] != expectedNtpdateArgs[1] {
				t.Fatalf("expected %s, got %s", expectedNtpdateArgs[1], arg[1])
			}

			return nil
		},
	}
	ntpdate.Ntpdate(expectedNtpdateArgs[1])
}
