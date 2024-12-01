package config

import (
	"flag"
	"fmt"
	"os"
)

type settings struct {
	AddressRun  string
	AddressBase string
}

var ServerSettings settings

func InitSettings() {
	ServerSettings = settings{}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Version: 1\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&ServerSettings.AddressRun, "a", "localhost:8080", "Run address")
	flag.StringVar(&ServerSettings.AddressBase, "b", "http://localhost:8080", "Base Address")
	flag.Parse()

	if envAddressRun := os.Getenv("SERVER_ADDRESS"); envAddressRun != "" {
		ServerSettings.AddressRun = envAddressRun
	}
	if envAddressBase := os.Getenv("BASE_URL"); envAddressBase != "" {
		ServerSettings.AddressBase = envAddressBase
	}
}
