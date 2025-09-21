package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Fund struct {
	Name   string  
	Amount string
}

type Fund2 struct {
	timestamp time.Time
	Name string
	Amount string
	Source string
	Notes string
}

var LastFunds []Fund
var LastFundsBCA []Fund
var LastFundsDANA []Fund
var LastFundsGOPAY []Fund
var LastFundsTUNAI []Fund

func main() {
	// Spreadsheet ID (ambil dari URL Google Sheets)
	sheetID := "18m_6ULVBDhmz-4tES_ijLeN2IpZiM9SKGYV02drRv5E"

	// Daftar gid sheet (cek URL #gid=xxx)
	gids := []string{"0", "1227436461"}
	names := []string{"Sheet1.csv", "Sheet2.csv"}

	for i, gid := range gids {
		url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", sheetID, gid)

		// HTTP GET
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			panic(fmt.Sprintf("Failed to download sheet (gid=%s): %s", gid, resp.Status))
		}

		// Simpan ke file
		out, err := os.Create(names[i])
		if err != nil {
			panic(err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Downloaded:", names[i])
	}
	readCsv("Sheet1.csv")
	readCsv("Sheet2.csv")
	fmt.Println(len(LastFunds))
}

func readCsv(path string) {
	file, _ := os.Open(path)
	reader := csv.NewReader(file)
	i := 0
	for {
		if (i < 3) {
			i++
			continue
		}
		record, err := reader.Read()
		if (err != nil) {
			if (err.Error() == "EOF") {
				break
			}
			log.Fatal(err)
		}
		LastFunds = append(LastFunds, mapToFund(record))
		fmt.Println(LastFunds[0])
		i++
	}
}

func mapToFund(strs []string) Fund {
	return Fund {
		Name: strs[1],
		Amount: strs[2],
	}
}
