package main

import "log"

func getSecretLocal(url string) (string, error) {
	log.Println("get secret local")
	return "this is local secret !!", nil
}
