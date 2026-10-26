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
	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		fmt.Println("Fund Raiser Tracker available flags:")
		flag.PrintDefaults()
		return
	}

	if *config.LoadStartIndexFlag {
		index, err := config.LoadStartIndex(*config.StateFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(index)
		return
	}
	if *config.SaveStartIndexFlag >= 0 {
		if err := config.SaveStartIndex(*config.SaveStartIndexFlag, *config.StateFile); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("saved start index: %d\n", *config.SaveStartIndexFlag)
		return
	}

	log.Println("SheetID: " + *config.SheetId)
	log.Println("GID: " + *config.Gid)
	log.Println("IntervalSeconds: " + fmt.Sprint(*config.IntervalSeconds))
	log.Println("TotalFunds: " + fmt.Sprint(*config.TotalFunds))
	log.Println("StartIndex: " + fmt.Sprint(*config.StartIndex))
	log.Println("OutputFormat: " + *config.OutputFormat)
	runProcess()
}

func runProcess() {
	log.Println("Starting fund tracker...")
	fund.WriteFunds()
	if *config.RunOnce == "y" {
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
		fund.WriteFunds()
	}
}
