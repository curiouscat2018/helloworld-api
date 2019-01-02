package config

import "os"

const Host = "api.curiouscat.one"
const DemosecretVaultIdentifier = "https://helloworld-vault.vault.azure.net/secrets/demosecret/ff32ea6957d04e529407cc839eef2cf8?api-version=7.0"
const DBURLVaultIdentifier = "https://helloworld-vault.vault.azure.net/secrets/helloworld-db-connection-str/0e26fb007a3c49caabfad16be2e6713e?api-version=7.0"

type configuration struct {
}

var Config configuration

func (c configuration) IsMockEnv() bool {
	return os.Getenv("HELLOWORLD_API_ENV") != "PROD"
}

func (c configuration) GetHomeDir() string {
	return os.Getenv("HOME")
}
