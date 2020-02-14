package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-aircraft"
	"github.com/sfomuseum/go-sfomuseum-aircraft/icao"
	"github.com/sfomuseum/go-sfomuseum-aircraft/sfomuseum"
	"log"
)

func main() {

	source := flag.String("source", "sfomuseum", "Valid sources are: icao, sfomuseum.")
	flag.Parse()

	var lookup aircraft.Lookup
	var err error

	switch *source {
	case "icao":
		lookup, err = icao.NewLookup()
	case "sfomuseum":
		lookup, err = sfomuseum.NewLookup()
	default:
		err = errors.New("Unknown source")
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range flag.Args() {

		results, err := lookup.Find(code)

		if err != nil {
			log.Fatal(err)
		}

		for _, a := range results {
			fmt.Println(a)
		}
	}
}
