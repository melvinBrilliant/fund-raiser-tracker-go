package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	SheetId            *string
	Gid                *string
	IntervalSeconds    *int
	TotalFunds         *int
	StartIndex         *int
	StateFile          *string
	OutputFormat       *string
	HTMLTemplateFile   *string
	HTMLOutputFile     *string
	RunOnce            *string
	LoadStartIndexFlag *bool
	SaveStartIndexFlag *int
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
	SheetId = flag.String("sheetid", "1s5nSkQX8TXh1x7AWi91s6i6feOofQaYPz3YdjlGUp8s", "Google Sheet ID")
	Gid = flag.String("gid", "1492535891", "Google Sheet GID")
	IntervalSeconds = flag.Int("interval", 20, "Interval in seconds")
	TotalFunds = flag.Int("total", 4, "Total funds to show")
	StateFile = flag.String("state-file", "fund_tracker_state.json", "Path to persisted start index state file")
	OutputFormat = flag.String("output", "txt", "Output format: txt or html")
	HTMLTemplateFile = flag.String("html-template", "templates/funds.html", "HTML template file")
	HTMLOutputFile = flag.String("html-output", "funds.html", "HTML output file")
	RunOnce = flag.String("once", "n", "Run once without scheduling (use y)")
	StartIndex = flag.Int("start", -1, "Start index for donor window; defaults to the value stored in state file")
	LoadStartIndexFlag = flag.Bool("load-start-index", false, "Load the persisted start index and print it, then exit")
	SaveStartIndexFlag = flag.Int("save-start-index", -1, "Persist a start index to the state file and exit")
	flag.Parse()

	if *SaveStartIndexFlag >= 0 {
		return
	}
	if *StartIndex < 0 {
		loadedIndex, err := LoadStartIndex(*StateFile)
		if err != nil {
			panic(err)
		}
		*StartIndex = loadedIndex
	}
}
