package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-aircraft/icao"
	_ "github.com/sfomuseum/go-sfomuseum-aircraft/sfomuseum"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-aircraft"
	"log"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "sfomuseum://", "Valid options are: icao://, sfomuseum://")

	flag.Parse()

	ctx := context.Background()
	lookup, err := aircraft.NewLookup(ctx, *lookup_uri)

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range flag.Args() {

		results, err := lookup.Find(ctx, code)

		if err != nil {
			log.Fatal(err)
		}

		for _, a := range results {
			fmt.Println(a)
		}
	}
}
