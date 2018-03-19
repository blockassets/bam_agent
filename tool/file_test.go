package tool

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIsExecutable(t *testing.T) {
	if IsExecutable(0644) {
		t.Fatalf("expected not executable for 0644")
	}
	if !IsExecutable(0544) {
		t.Fatalf("expected executable for 0544")
	}
	if !IsExecutable(0777) {
		t.Fatalf("expected executable for 0777")
	}
}

func TestIsDirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "tool-test")
	defer os.RemoveAll(tempDir)

	if err != nil {
		t.Fatal(err)
	}

	isDir, err := IsDirectory(tempDir)
	if !isDir || err != nil {
		t.Fatal(err)
		t.Fatalf("expected directory, got: %v", isDir)
	}
}
