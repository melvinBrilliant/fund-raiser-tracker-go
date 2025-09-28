package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"tracker/config"
	"tracker/fund"
)

func main() {
	config.InitConfig()
	log.Println("SheetID: " + *config.SheetId)
	log.Println("GID: " + *config.Gid)
	log.Println("IntervalSeconds: " + fmt.Sprint(*config.IntervalSeconds))
	log.Println("TotalFunds: " + fmt.Sprint(*config.TotalFunds))
	runProcess()
}

func runProcess() {
	log.Println("Starting fund tracker...")
	fund.WriteFundsToTxt()
	runOnce := flag.String("once", "n", "Run once without scheduling")
	flag.Parse()
	if (runOnce != nil && *runOnce == "y") {
		log.Println("Running once without scheduling...")
		return
	} else {
		scheduledProcess()
	}
}

func scheduledProcess() {
	ticker := time.NewTicker(time.Duration(*config.IntervalSeconds) * time.Second)
	defer ticker.Stop()
	for t := range ticker.C {
		log.Printf("Scheduled task running at %v...\n", t)
		fund.WriteFundsToTxt()
	}
}
