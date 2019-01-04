package config

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Host        string
	DB_URL      string
	IsMockEnv   bool
	TLSCertPath string
}

var Config config

func init() {
	f, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewDecoder(f).Decode(&Config); err != nil {
		log.Fatal(err)
	}
}