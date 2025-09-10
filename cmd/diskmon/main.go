package main

import (
	"context"
	"diskmon/internal/config"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()
	
	cnf, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	
	err = cnf.Watch(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}	
}