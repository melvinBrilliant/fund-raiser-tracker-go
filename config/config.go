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
	SheetId = flag.String("sheetid", "1zHfOOSsIrXjlPqiaJgRfC5xzo39ZD56AUoCIyojhonI", "Google Sheet ID")
	Gid = flag.String("gid", "1039232989", "Google Sheet GID")
	IntervalSeconds = flag.Int("interval", 30, "Interval in seconds")
	TotalFunds = flag.Int("total", 5, "Total funds to show")
	flag.Parse()
}
