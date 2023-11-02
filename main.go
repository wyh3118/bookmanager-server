package main

import "bookmanager-server/router"
import "bookmanager-server/initialize"

func init() {
	initialize.Init()
}

func main() {
	engin := router.GetEngine()

	if err := engin.Run(":7001"); err != nil {
		panic(err)
	}
}
