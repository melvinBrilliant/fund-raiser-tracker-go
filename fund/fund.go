package fund

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/message"
	"tracker/config"
)

type Fund struct {
	timestamp time.Time
	Name      string
	Amount    float64
	Source    string
	Notes     string
}

func downloadCSV() {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", *config.SheetId, *config.Gid)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to download sheet (gid=%s): %s", *config.Gid, resp.Status)
	}
	defer resp.Body.Close()
	out, err := os.Create("Sheet.csv")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

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
	amountFloat, _ := strconv.ParseFloat(row[4], 64)
	return Fund{
		timestamp: mapTimestamp(row[0]),
		Name:      row[2],
		Amount:    amountFloat,
		Source:    row[1],
		Notes:     row[3],
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

func latestFunds() []Fund {
	funds := mapCsv()
	if len(funds) < *config.TotalFunds {
		log.Fatalf("Not enough funds in CSV (found=%d, required=%d)", len(funds), config.TotalFunds)
	}
	return funds[len(funds)-*config.TotalFunds:]
}

func fundsText() string {
	funds := latestFunds()
	template := "%s %s"
	var text []string
	for _, fund := range funds {
		rupiahAmount := message.NewPrinter(message.MatchLanguage("id")).Sprintf("Rp%.0f", fund.Amount)
		fundText := fmt.Sprintf(template, fund.Name, rupiahAmount)
		text = append(text, fundText)
		log.Printf("Added fund text: %s", fundText)
	}
	return strings.Join(text, " | ")
}

func WriteFundsToTxt() {
	downloadCSV()
	text := fundsText()
	os.WriteFile("funds.txt", []byte(text), 0644)
}
