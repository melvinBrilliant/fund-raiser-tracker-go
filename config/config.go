package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	SheetId            *string
	Gid                *string
	IntervalSeconds    *int
	TotalFunds         *int
	StartIndex         *int
	StateFile          *string
	LoadStartIndexFlag *bool
	SaveStartIndexFlag *int
	RunOnce            *bool
)

type persistedStartIndex struct {
	StartIndex int `json:"start_index"`
}

func LoadStartIndex(statePath string) (int, error) {
	if statePath == "" {
		statePath = "fund_tracker_state.json"
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, fmt.Errorf("read state file: %w", err)
	}
	if len(data) == 0 {
		return 0, nil
	}

	var state persistedStartIndex
	if err := json.Unmarshal(data, &state); err != nil {
		return 0, fmt.Errorf("parse state file: %w", err)
	}
	if state.StartIndex < 0 {
		return 0, nil
	}
	return state.StartIndex, nil
}

func SaveStartIndex(index int, statePath string) error {
	if statePath == "" {
		statePath = "fund_tracker_state.json"
	}
	if index < 0 {
		index = 0
	}
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

	state := persistedStartIndex{StartIndex: index}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	tempPath := statePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp state: %w", err)
	}
	if err := os.Rename(tempPath, statePath); err != nil {
		return fmt.Errorf("replace state file: %w", err)
	}
	return nil
}

func InitConfig() {
	args := os.Args[1:]
	for i, arg := range args {
		if strings.HasPrefix(arg, "-once=") {
			value := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(arg, "-once=")))
			if value == "y" || value == "yes" || value == "true" || value == "1" {
				args[i] = "-once=true"
			} else if value == "n" || value == "no" || value == "false" || value == "0" {
				args[i] = "-once=false"
			}
		}
	}
	os.Args = append([]string{os.Args[0]}, args...)

	SheetId = flag.String("sheetid", "1s5nSkQX8TXh1x7AWi91s6i6feOofQaYPz3YdjlGUp8s", "Google Sheet ID")
	Gid = flag.String("gid", "1492535891", "Google Sheet GID")
	IntervalSeconds = flag.Int("interval", 1	git add config/config.go main.go config/config_test.go
	git commit -m "fix: initialize state file and normalize once flag"5, "Interval in seconds")
	TotalFunds = flag.Int("total", 3, "Total funds to show")
	StateFile = flag.String("state-file", "fund_tracker_state.json", "Path to persisted start index state file")
	StartIndex = flag.Int("start", -1, "Start index for donor window; defaults to the value stored in state file")
	LoadStartIndexFlag = flag.Bool("load-start-index", false, "Load the persisted start index and print it, then exit")
	SaveStartIndexFlag = flag.Int("save-start-index", -1, "Persist a start index to the state file and exit")
	RunOnce = flag.Bool("once", false, "Run once without scheduling")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	if *SaveStartIndexFlag >= 0 {
		return
	}
	if *StartIndex < 0 {
		loadedIndex, err := LoadStartIndex(*StateFile)
		if err != nil {
			panic(err)
		}
		if loadedIndex == 0 {
			if _, statErr := os.Stat(*StateFile); errors.Is(statErr, os.ErrNotExist) {
				if err := SaveStartIndex(0, *StateFile); err != nil {
					panic(err)
				}
			}
		}
		*StartIndex = loadedIndex
	}
}
