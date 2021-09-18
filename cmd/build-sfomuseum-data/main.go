package main

import (
	"context"
	"encoding/json"
	"flag"
	tools "github.com/sfomuseum/go-sfomuseum-aircraft-tools"
	"github.com/sfomuseum/go-sfomuseum-aircraft-tools/template"
	"log"
	"os"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://", "...")
	iterator_source := flag.String("iterator-source", "/usr/local/data/sfomuseum-data-aircraft", "...")

	flag.Parse()

	ctx := context.Background()

	lookup, err := tools.CompileAircraftData(ctx, *iterator_uri, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to compile aircraft data, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(lookup)

	if err != nil {
		log.Fatalf("Failed to marshal results, %w", err)
	}
}
