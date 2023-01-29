package main

import (
	"flag"
	"log"

	"github.com/farischt/gobank/config"
)

func init() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		log.Fatalf("Usage: server -e {mode}")
	}
	flag.Parse()
	config.Init(*environment)
}

func main() {
	port := ":" + config.GetConfig().GetString("PORT")
	storage, err := NewPgStorage()

	if err != nil {
		log.Fatal(err)
	}

	s := NewApiServer(port, storage)
	s.Start()
}
