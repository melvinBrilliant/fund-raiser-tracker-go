package config

import (
	"flag"
)

var (
	SheetId         *string
	Gid             *string
	IntervalSeconds *int
	TotalFunds      *int
)

func InitConfig() {
	SheetId = flag.String("sheetid", "1GHJojMdDrp6pyORdZT5o1_0Njr0VTlG294IMtBcE124", "Google Sheet ID")
	Gid = flag.String("gid", "718286749", "Google Sheet GID")
	IntervalSeconds = flag.Int("interval", 20, "Interval in seconds")
	TotalFunds = flag.Int("total", 4, "Total funds to show")
	flag.Parse()
}
