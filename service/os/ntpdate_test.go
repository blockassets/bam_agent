package os

import "testing"

func TestNtpdate_Ntpdate(t *testing.T) {
	expectedNtpdateCmd := "/usr/bin/ntpdate"
	expectedNtpdateArg := "-u time.google.com"

	ntpdate := &NtpdateData{
		run: func(cmd string, arg string) error {
			if cmd != expectedNtpdateCmd {
				t.Fatalf("expected %s, got %s", expectedNtpdateCmd, cmd)
			}
			if arg != expectedNtpdateArg {
				t.Fatalf("expected %s, got %s", expectedNtpdateArg, arg)
			}

			return nil
		},
	}
	ntpdate.Ntpdate()
}
