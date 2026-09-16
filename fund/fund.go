package fund

import (
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"tracker/config"
)

type Fund struct {
	timestamp time.Time
	Name      string
	Amount    string
}

func downloadCSV() error {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", *config.SheetId, *config.Gid)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download sheet (gid=%s): %s", *config.Gid, resp.Status)
		return fmt.Errorf("FAILED TO DOWNLOAD SHEET (GID=%s): %s", *config.Gid, resp.Status)
	}
	defer resp.Body.Close()
	out, err := os.Create("Sheet.csv")
	if err != nil {
		log.Printf("Failed to create Sheet.csv: %s", err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Failed to write to Sheet.csv: %s", err)
		return err
	}
	return nil
}

func mapTimestamp(s string) time.Time {
	layout := "02/01/2006 15:04:05"
	t, err := time.Parse(layout, s)
	if err != nil {
		log.Fatalf("Failed to parse timestamp: %s", err)
	}
	return t
}

func mapCsvRow(row []string) Fund {
	return Fund{
		timestamp: mapTimestamp(row[0]),
		Name:      row[1],
		Amount:    row[2],
	}
}

func mapCsv() []Fund {
	file, _ := os.Open("Sheet.csv")
	defer file.Close()
	reader := csv.NewReader(file)
	var funds []Fund
	for i := 0; ; i++ {
		if i == 0 {
			// skip header
			_, err := reader.Read()
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		funds = append(funds, mapCsvRow(record))
	}
	return funds
}

func selectFundsWindow(funds []Fund, startIndex, total int) ([]Fund, int) {
	if len(funds) == 0 {
		return nil, 0
	}
	if total <= 0 {
		return nil, startIndex
	}
	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex >= len(funds) {
		startIndex = 0
	}

	selectedEnd := min(startIndex+total, len(funds))

	selected := funds[startIndex:selectedEnd]
	if len(selected) < total {
		return selected, startIndex
	}

	nextIndex := startIndex + 1
	if nextIndex+total > len(funds) {
		return selected, startIndex
	}

	return selected, nextIndex
}

func extractFunds() ([]Fund, error) {
	funds := mapCsv()
	if len(funds) == 0 {
		errorMsg := "No funds found in CSV"
		log.Print(errorMsg)
		return nil, errors.New(errorMsg)
	}
	if len(funds) < *config.TotalFunds {
		errorMsg := fmt.Sprintf("Not enough funds in CSV (found=%d, required=%d)", len(funds), *config.TotalFunds)
		log.Print(errorMsg)
		return nil, errors.New(errorMsg)
	}

	selected, nextIndex := selectFundsWindow(funds, *config.StartIndex, *config.TotalFunds)
	*config.StartIndex = nextIndex
	if err := config.SaveStartIndex(*config.StartIndex, *config.StateFile); err != nil {
		log.Printf("failed to persist start index: %v", err)
	}
	return selected, nil
}

func fundsText() (string, error) {
	funds, err := extractFunds()
	if err != nil {
		return "", err
	}
	return fundsTextFrom(funds), nil
}

func formatSpaces(s string) string {
	maxLength := 135
	if len(s) >= maxLength {
		return s
	}
	remainingSpaces := maxLength - len(s)
	for i := 0; i < remainingSpaces; i++ {
		s = s + " "
	}
	return s
}

func writeFundsToTxt(text string) error {
	return os.WriteFile("funds.txt", []byte(formatSpaces(text)), 0644)
}

type htmlTemplateData struct {
	Funds             []Fund
	AnimationDuration int
	BackgroundColor   string
	FontColor         string
}

func WriteFundsToHTML(funds []Fund) error {
	if *config.TickerSpeed <= 0 {
		return fmt.Errorf("ticker speed must be greater than zero")
	}
	tmpl, err := template.ParseFiles(*config.HTMLTemplateFile)
	if err != nil {
		return fmt.Errorf("parse HTML template: %w", err)
	}
	for _, fund := range funds {
		log.Printf("Added fund text: %s %s", fund.Name, fund.Amount)
	}

	output, err := os.Create(*config.HTMLOutputFile)
	if err != nil {
		return fmt.Errorf("create HTML output: %w", err)
	}
	defer output.Close()

	data := htmlTemplateData{
		Funds:             funds,
		AnimationDuration: *config.TickerSpeed,
		BackgroundColor:   *config.BackgroundColor,
		FontColor:         *config.FontColor,
	}
	if err := tmpl.Execute(output, data); err != nil {
		return fmt.Errorf("render HTML template: %w", err)
	}
	return nil
}

func WriteFunds() {
	err := downloadCSV()
	if err != nil {
		log.Printf("Error downloading CSV: %s", err)
		return
	}
	funds, err := extractFunds()
	if err != nil {
		log.Printf("Error extracting funds: %s", err)
		return
	}

	if *config.OutputFormat == "html" {
		if err := WriteFundsToHTML(funds); err != nil {
			log.Printf("Error generating funds HTML: %s", err)
		}
		return
	}

	if err := writeFundsToTxt(fundsTextFrom(funds)); err != nil {
		log.Printf("Error writing funds.txt: %s", err)
	}
}

func fundsTextFrom(funds []Fund) string {
	var text []string
	for _, fund := range funds {
		fundText := fmt.Sprintf("%s %s", fund.Name, fund.Amount)
		text = append(text, fundText)
		log.Printf("Added fund text: %s", fundText)
	}
	return strings.Join(text, " | ") + " | "
}

func WriteFundsToTxt() {
	text, err := fundsText()
	if err != nil {
		log.Printf("Error generating funds text: %s", err)
		return
	}
	if err := writeFundsToTxt(text); err != nil {
		log.Printf("Error writing funds.txt: %s", err)
	}
}
