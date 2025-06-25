package main

import (
	"github.com/Matthew-Mak/go-moni-shark/bot"
	"github.com/Matthew-Mak/go-moni-shark/config"
	"log"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Start()

	<-make(chan struct{})

	return
}
