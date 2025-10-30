package fund

import (
	"encoding/csv"
	"errors"
	"fmt"
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
	Source    string
	Notes     string
}

func downloadCSV() error {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", *config.SheetId, *config.Gid)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download sheet (gid=%s): %s", *config.Gid, resp.Status)
		return fmt.Errorf("FAILED TO DOWNLOAD SHEET (GID=%s): %s", *config.Gid, resp.Status);
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
		Amount:    row[3],
		Source:    row[2],
		Notes:     row[4],
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

func latestFunds() ([]Fund, error) {
	funds := mapCsv()
	if len(funds) < *config.TotalFunds {
		errorMsg := fmt.Sprintf("Not enough funds in CSV (found=%d, required=%d)", len(funds), config.TotalFunds)
		log.Print(errorMsg)
		return nil, errors.New(errorMsg)
	}
	return funds[len(funds)-*config.TotalFunds:], nil
}

func fundsText() (string, error) {
	funds, err := latestFunds()
	if err != nil {	
		return "", err
	}
	template := "%s %s"
	var text []string
	for _, fund := range funds {
		// rupiahAmount := message.NewPrinter(message.MatchLanguage("id")).Sprintf("Rp%.0f", fund.Amount)
		fundText := fmt.Sprintf(template, fund.Name, fund.Amount)
		text = append(text, fundText)
		log.Printf("Added fund text: %s", fundText)
	}
	return strings.Join(text, " | ") + " | ", nil
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

func WriteFundsToTxt() {
	err := downloadCSV()
	if err != nil {
		log.Printf("Error downloading CSV: %s", err)
		return
	}
	text, err := fundsText()
	text = formatSpaces(text)
	if err != nil {
		log.Printf("Error generating funds text: %s", err)
		return
	}
	os.WriteFile("funds.txt", []byte(text), 0644)
}
