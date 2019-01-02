package config

import "os"

const Host = "api.curiouscat.one"
const DBURL = "mongodb://helloworld-db:mMIXPtgqLRa8FWhIzmbuKWTNvSyL2kmdbdewIton3iFp9lqimEhofbMTlQNcNNiSdtmZBfiVpGau5OVLHqPLNg==@helloworld-db.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
const TestSecretUrl = "https://helloworld-vault.vault.azure.net/secrets/testsecret/e3246fd47fa74a638e99e4c4afe97006?api-version=7.0"
const DemoSecretUrl = "https://helloworld-vault.vault.azure.net/secrets/demosecret/ff32ea6957d04e529407cc839eef2cf8?api-version=7.0"

type configuration struct {
	i int
}

var Config configuration

func (c configuration) IsMockEnv() bool {
	return os.Getenv("HELLOWORLD_API_ENV") != "PROD"
}

func (c configuration) GetHomeDir() string {
	return os.Getenv("HOME")
}
