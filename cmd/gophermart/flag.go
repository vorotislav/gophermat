package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"gophermat/internal/settings"
)

func parseFlag(set *settings.Settings) {
	var address string
	flag.StringVar(&address, "a", ":8081", "address and port to run server")
	var databaseURI string
	flag.StringVar(&databaseURI, "d", "", "database uri")
	var accrualAddress string
	flag.StringVar(&accrualAddress, "r", ":8080", "address and port of accrual system")

	flag.Parse()

	if err := env.Parse(set); err == nil {
		if set.Address == "" {
			set.Address = address
		}

		if set.DatabaseURI == "" {
			set.DatabaseURI = databaseURI
		}

		if set.AccrualSystemAddress == "" {
			set.AccrualSystemAddress = accrualAddress
		}
	}
}
