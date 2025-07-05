package main

import (
	"diskmon/internal/config"
	"log"
	"os"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	
	_=c
}