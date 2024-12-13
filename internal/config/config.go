package config

import (
	"flag"
	"os"
)

type ServerSettings struct {
	AddressRun  string
	AddressBase string
}

func GetConfig() ServerSettings {
	settings := ServerSettings{}

	flag.StringVar(&settings.AddressRun, "a", "localhost:8080", "Run address")
	flag.StringVar(&settings.AddressBase, "b", "http://localhost:8080", "Base Address")
	flag.Parse()

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok && val != "" {
		settings.AddressRun = val
	}

	val, ok = os.LookupEnv("BASE_URL")
	if ok && val != "" {
		settings.AddressBase = val
	}

	return settings
}
