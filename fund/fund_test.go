package fund

import (
	"os"
	"strings"
	"testing"

	"tracker/config"
)

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

func TestWriteFundsToHTMLRendersFunds(t *testing.T) {
	templatePath := t.TempDir() + "/funds.html"
	outputPath := t.TempDir() + "/output.html"
	if err := os.WriteFile(templatePath, []byte("{{range .}}<p>{{.Name}} {{.Amount}}</p>{{end}}"), 0644); err != nil {
		t.Fatal(err)
	}

	previousTemplate := config.HTMLTemplateFile
	previousOutput := config.HTMLOutputFile
	templateFile := templatePath
	outputFile := outputPath
	config.HTMLTemplateFile = &templateFile
	config.HTMLOutputFile = &outputFile
	t.Cleanup(func() {
		config.HTMLTemplateFile = previousTemplate
		config.HTMLOutputFile = previousOutput
	})

	funds := []Fund{{Name: "Melvin", Amount: "Rp100.000"}}
	if err := WriteFundsToHTML(funds); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Melvin Rp100.000") {
		t.Fatalf("rendered HTML does not contain fund: %s", data)
	}
}
