package config

import (
	"encoding/json"
	"os"

	"github.com/curiouscat2018/helloworld-api/common"
)

type config struct {
	DB_URL    string
	IsMockEnv bool
}

var Config = config{
	DB_URL:    "",
	IsMockEnv: true,
}

func init() {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		common.TraceInfo(nil).Msg("config file not found. use default values")
		return
	}

	f, err := os.Open("./config.json")
	if err != nil {
		common.TraceFatal(nil).Err(err).Msg("unable to open config file")
	}

	if err := json.NewDecoder(f).Decode(&Config); err != nil {
		common.TraceFatal(nil).Err(err).Msg("unable to parse config file")
	}
}

func (c config) HostName() string {
	res, _ := os.Hostname()
	if res == "" {
		res = "NA"
	}

	return res
}
