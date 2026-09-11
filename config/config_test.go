package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfigCreatesDefaultStateFileWhenMissing(t *testing.T) {
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	os.Args = []string{"cmd", "-state-file", statePath}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	InitConfig()

	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("expected state file to be created on first run: %v", err)
	}
}
