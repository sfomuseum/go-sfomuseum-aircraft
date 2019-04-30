package main

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-aircraft/icao"
	"log"
)

func main() {

	flag.Parse()

	l, err := icao.NewLookup()

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range flag.Args() {

		aircraft, err := l.Find(code)

		if err != nil {
			log.Fatal(err)
		}

		for _, a := range aircraft {
			fmt.Println(a)
		}
	}
}
