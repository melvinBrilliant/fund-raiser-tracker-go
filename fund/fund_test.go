package fund

import "testing"

func TestWindowFundsAdvancesWhenEnoughDataRemains(t *testing.T) {
	funds := make([]Fund, 10)
	for i := range funds {
		funds[i] = Fund{Name: string(rune('A' + i))}
	}

	selected, nextIndex := selectFundsWindow(funds, 0, 4)
	if len(selected) != 4 {
		t.Fatalf("expected 4 funds, got %d", len(selected))
	}
	if selected[0].Name != "A" || selected[3].Name != "D" {
		t.Fatalf("unexpected first window: %#v", selected)
	}
	if nextIndex != 1 {
		t.Fatalf("expected next index 1, got %d", nextIndex)
	}
}

func TestWindowFundsDoesNotAdvanceWhenNotEnoughDataRemains(t *testing.T) {
	funds := make([]Fund, 5)
	for i := range funds {
		funds[i] = Fund{Name: string(rune('A' + i))}
	}

	selected, nextIndex := selectFundsWindow(funds, 1, 4)
	if len(selected) != 4 {
		t.Fatalf("expected 4 funds, got %d", len(selected))
	}
	if selected[0].Name != "B" || selected[3].Name != "E" {
		t.Fatalf("unexpected selection when not enough remaining: %#v", selected)
	}
	if nextIndex != 1 {
		t.Fatalf("expected next index to remain 1, got %d", nextIndex)
	}
}

func TestWindowFundsResetsWhenIndexRunsPastEnd(t *testing.T) {
	funds := make([]Fund, 5)
	for i := range funds {
		funds[i] = Fund{Name: string(rune('A' + i))}
	}

	selected, nextIndex := selectFundsWindow(funds, 10, 4)
	if len(selected) != 4 {
		t.Fatalf("expected 4 funds, got %d", len(selected))
	}
	if selected[0].Name != "A" || selected[3].Name != "D" {
		t.Fatalf("unexpected reset selection: %#v", selected)
	}
	if nextIndex != 1 {
		t.Fatalf("expected reset next index 1, got %d", nextIndex)
	}
}
