package config

import (
	"flag"
	"os"
)

type ServerSettings struct {
	AddressRun  string
	AddressBase string
	LogLevel    string
	FileStorage string
	DBStorage   string
}

func GetConfig() ServerSettings {
	settings := ServerSettings{}

	flag.StringVar(&settings.AddressRun, "a", "localhost:8080", "Run address")
	flag.StringVar(&settings.AddressBase, "b", "http://localhost:8080", "Base Address")
	flag.StringVar(&settings.LogLevel, "l", "info", "Log level")
	flag.StringVar(&settings.FileStorage, "f", "/tmp/storage.json", "File storage data")
	flag.StringVar(&settings.DBStorage, "d", "", "Address connect to DB")
	flag.Parse()

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok && val != "" {
		settings.AddressRun = val
	}

	val, ok = os.LookupEnv("BASE_URL")
	if ok && val != "" {
		settings.AddressBase = val
	}

	val, ok = os.LookupEnv("LOG_LEVEL")
	if ok && val != "" {
		settings.LogLevel = val
	}

	val, ok = os.LookupEnv("FILE_STORAGE_PATH")
	if ok && val != "" {
		settings.FileStorage = val
	}

	val, ok = os.LookupEnv("DATABASE_DSN")
	if ok && val != "" {
		settings.DBStorage = val
	}

	return settings
}
