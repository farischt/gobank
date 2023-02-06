package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/farischt/gobank/api"
	"github.com/farischt/gobank/config"
	"github.com/farischt/gobank/store"
)

func init() {
	environment := flag.String("e", "dev", "")
	flag.Usage = func() {
		log.Fatalf("Usage: server -e {mode}")
	}
	flag.Parse()
	config.InitBaseConfig(*environment)
	config.InitDbConfig(*environment)
}

func main() {
	configPort := config.GetConfig().GetInt(config.PORT)
	port := fmt.Sprintf(":%d", configPort)
	storage, err := store.NewPgStore()

	if err != nil {
		log.Fatal(err)
	}

	s := api.New(port, *storage)
	s.Start()
}
