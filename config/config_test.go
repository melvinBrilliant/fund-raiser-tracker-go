package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadStartIndexCreatesMissingStateFile(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "nested", "state.json")

	index, err := LoadStartIndex(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if index != 0 {
		t.Fatalf("expected start index 0, got %d", index)
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatal(err)
	}
	if state.StartIndex != 0 {
		t.Fatalf("expected persisted start index 0, got %d", state.StartIndex)
	}
	if state.TickerSpeed != DefaultTickerSpeed {
		t.Fatalf("expected persisted ticker speed %d, got %d", DefaultTickerSpeed, state.TickerSpeed)
	}
	if state.SheetID != DefaultSheetID || state.GID != DefaultGID {
		t.Fatalf("expected default sheet settings, got %#v", state)
	}
}

func TestLoadStateUsesPersistedTickerSpeed(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	if err := SaveStateWithInterval("sheet-123", "456", 8, 12, 20, 48, "#123456", "#abcdef", statePath); err != nil {
		t.Fatal(err)
	}

	state, err := LoadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if state.SheetID != "sheet-123" || state.GID != "456" || state.StartIndex != 8 || state.IntervalSeconds != 12 || state.TickerSpeed != 20 || state.FontSize != 48 || state.BackgroundColor != "#123456" || state.FontColor != "#abcdef" {
		t.Fatalf("unexpected state: %#v", state)
	}
}
