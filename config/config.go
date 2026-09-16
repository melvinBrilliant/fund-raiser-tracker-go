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
	TickerSpeed        *int
	FontSize           *int
	BackgroundColor    *string
	FontColor          *string
	RunOnce            *string
	LoadStartIndexFlag *bool
	SaveStartIndexFlag *int
)

const DefaultSheetID = "1s5nSkQX8TXh1x7AWi91s6i6feOofQaYPz3YdjlGUp8s"
const DefaultGID = "1492535891"
const DefaultIntervalSeconds = 20
const DefaultTickerSpeed = 45
const DefaultFontSize = 48
const DefaultBackgroundColor = "#011D66"
const DefaultFontColor = "#f2cc8f"

type persistedState struct {
	SheetID         string `json:"sheet_id"`
	GID             string `json:"gid"`
	StartIndex      int    `json:"start_index"`
	IntervalSeconds int    `json:"interval_seconds"`
	TickerSpeed     int    `json:"ticker_speed"`
	FontSize        int    `json:"font_size"`
	BackgroundColor string `json:"background_color"`
	FontColor       string `json:"font_color"`
}

func LoadState(statePath string) (persistedState, error) {
	if statePath == "" {
		statePath = "fund_tracker_state.json"
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			state := defaultState()
			if err := saveState(state, statePath); err != nil {
				return persistedState{}, fmt.Errorf("create state file: %w", err)
			}
			return state, nil
		}
		return persistedState{}, fmt.Errorf("read state file: %w", err)
	}
	if len(data) == 0 {
		state := defaultState()
		if err := saveState(state, statePath); err != nil {
			return persistedState{}, fmt.Errorf("initialize state file: %w", err)
		}
		return state, nil
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return persistedState{}, fmt.Errorf("parse state file: %w", err)
	}
	needsSave := false
	if state.SheetID == "" {
		state.SheetID = DefaultSheetID
		needsSave = true
	}
	if state.GID == "" {
		state.GID = DefaultGID
		needsSave = true
	}
	if state.StartIndex < 0 {
		state.StartIndex = 0
		needsSave = true
	}
	if state.TickerSpeed <= 0 {
		state.TickerSpeed = DefaultTickerSpeed
		needsSave = true
	}
	if state.IntervalSeconds <= 0 {
		state.IntervalSeconds = DefaultIntervalSeconds
		needsSave = true
	}
	if state.FontSize <= 0 {
		state.FontSize = DefaultFontSize
		needsSave = true
	}
	if state.BackgroundColor == "" {
		state.BackgroundColor = DefaultBackgroundColor
		needsSave = true
	}
	if state.FontColor == "" {
		state.FontColor = DefaultFontColor
		needsSave = true
	}
	if needsSave {
		if err := saveState(state, statePath); err != nil {
			return persistedState{}, fmt.Errorf("upgrade state file: %w", err)
		}
	}
	return state, nil
}

func defaultState() persistedState {
	return persistedState{
		SheetID:         DefaultSheetID,
		GID:             DefaultGID,
		IntervalSeconds: DefaultIntervalSeconds,
		TickerSpeed:     DefaultTickerSpeed,
		FontSize:        DefaultFontSize,
		BackgroundColor: DefaultBackgroundColor,
		FontColor:       DefaultFontColor,
	}
}

func LoadStartIndex(statePath string) (int, error) {
	state, err := LoadState(statePath)
	return state.StartIndex, err
}

func SaveState(sheetID, gid string, index, tickerSpeed int, backgroundColor, fontColor, statePath string) error {
	return SaveStateWithInterval(sheetID, gid, index, DefaultIntervalSeconds, tickerSpeed, DefaultFontSize, backgroundColor, fontColor, statePath)
}

func SaveStateWithFontSize(sheetID, gid string, index, tickerSpeed, fontSize int, backgroundColor, fontColor, statePath string) error {
	return SaveStateWithInterval(sheetID, gid, index, DefaultIntervalSeconds, tickerSpeed, fontSize, backgroundColor, fontColor, statePath)
}

func SaveStateWithInterval(sheetID, gid string, index, intervalSeconds, tickerSpeed, fontSize int, backgroundColor, fontColor, statePath string) error {
	if statePath == "" {
		statePath = "fund_tracker_state.json"
	}
	if index < 0 {
		index = 0
	}
	if intervalSeconds <= 0 {
		return fmt.Errorf("interval must be greater than zero")
	}
	if tickerSpeed <= 0 {
		return fmt.Errorf("ticker speed must be greater than zero")
	}
	if fontSize <= 0 {
		return fmt.Errorf("font size must be greater than zero")
	}
	if backgroundColor == "" || fontColor == "" {
		return fmt.Errorf("background and font colors must not be empty")
	}
	if sheetID == "" || gid == "" {
		return fmt.Errorf("sheet ID and GID must not be empty")
	}

	return saveState(persistedState{
		SheetID:         sheetID,
		GID:             gid,
		StartIndex:      index,
		IntervalSeconds: intervalSeconds,
		TickerSpeed:     tickerSpeed,
		FontSize:        fontSize,
		BackgroundColor: backgroundColor,
		FontColor:       fontColor,
	}, statePath)
}

func SaveStartIndex(index int, statePath string) error {
	state, err := LoadState(statePath)
	if err != nil {
		return err
	}
	return SaveStateWithInterval(state.SheetID, state.GID, index, state.IntervalSeconds, state.TickerSpeed, state.FontSize, state.BackgroundColor, state.FontColor, statePath)
}

func saveState(state persistedState, statePath string) error {
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

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
	SheetId = flag.String("sheetid", DefaultSheetID, "Google Sheet ID")
	Gid = flag.String("gid", DefaultGID, "Google Sheet GID")
	IntervalSeconds = flag.Int("interval", DefaultIntervalSeconds, "Interval in seconds")
	TotalFunds = flag.Int("total", 4, "Total funds to show")
	StateFile = flag.String("state-file", "fund_tracker_state.json", "Path to persisted start index state file")
	OutputFormat = flag.String("output", "txt", "Output format: txt or html")
	HTMLTemplateFile = flag.String("html-template", "templates/funds-template.html", "HTML template file")
	HTMLOutputFile = flag.String("html-output", "funds.html", "HTML output file")
	TickerSpeed = flag.Int("ticker-speed", DefaultTickerSpeed, "Scrolling animation duration in seconds")
	FontSize = flag.Int("font-size", DefaultFontSize, "HTML ticker font size in pixels")
	BackgroundColor = flag.String("background-color", DefaultBackgroundColor, "HTML background color")
	FontColor = flag.String("font-color", DefaultFontColor, "HTML ticker font color")
	RunOnce = flag.String("once", "n", "Run once without scheduling (use y)")
	StartIndex = flag.Int("start", -1, "Start index for donor window; defaults to the value stored in state file")
	LoadStartIndexFlag = flag.Bool("load-start-index", false, "Load the persisted start index and print it, then exit")
	SaveStartIndexFlag = flag.Int("save-start-index", -1, "Persist a start index to the state file and exit")
	flag.Parse()

	if *SaveStartIndexFlag >= 0 {
		return
	}
	state, err := LoadState(*StateFile)
	if err != nil {
		panic(err)
	}
	if *StartIndex < 0 {
		*StartIndex = state.StartIndex
	}
	if !flagWasSet("interval") {
		*IntervalSeconds = state.IntervalSeconds
	}
	if !flagWasSet("ticker-speed") {
		*TickerSpeed = state.TickerSpeed
	}
	if !flagWasSet("font-size") {
		*FontSize = state.FontSize
	}
	if !flagWasSet("background-color") {
		*BackgroundColor = state.BackgroundColor
	}
	if !flagWasSet("font-color") {
		*FontColor = state.FontColor
	}
	if !flagWasSet("sheetid") {
		*SheetId = state.SheetID
	}
	if !flagWasSet("gid") {
		*Gid = state.GID
	}

	if *TickerSpeed <= 0 {
		panic("ticker speed must be greater than zero")
	}
	if *IntervalSeconds <= 0 {
		panic("interval must be greater than zero")
	}
	if *FontSize <= 0 {
		panic("font size must be greater than zero")
	}
	if err := SaveStateWithInterval(*SheetId, *Gid, *StartIndex, *IntervalSeconds, *TickerSpeed, *FontSize, *BackgroundColor, *FontColor, *StateFile); err != nil {
		panic(err)
	}
}

func flagWasSet(name string) bool {
	wasSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			wasSet = true
		}
	})
	return wasSet
}
