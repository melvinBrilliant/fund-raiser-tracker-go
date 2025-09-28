package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"id.melvin/config"
	"id.melvin/fund"
)

func main() {
	config.InitConfig()
	fmt.Println(*config.SheetId)
	runProcess()
}

func runProcess() {
	runOnce := flag.String("once", "n", "Run once without scheduling")
	flag.Parse()
	if (runOnce != nil && *runOnce == "y") {
		log.Println("Running once without scheduling...")
		fund.WriteFundsToTxt()
		return
	} else {
		scheduledProcess()
	}
}

func scheduledProcess() {
	ticker := time.NewTicker(time.Duration(*config.IntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("Scheduled task running...")
			fund.WriteFundsToTxt()
		}
	}
}
