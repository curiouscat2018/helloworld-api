package main

import "log"

func getSecretLocal(url string) (string, error) {
	log.Println("get local secret")
	return "this is local secret !!", nil
}
