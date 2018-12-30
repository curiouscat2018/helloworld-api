package main

import "os"

type configuration struct {
}

var config configuration

func (c configuration) isMockEnv() bool {
	return os.Getenv("HELLOWORLD_API_ENV") != "PROD"
}
