package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	// curl 'https://www4.icao.int/doc8643/External/AircraftTypes' -H 'Connection: keep-alive' --data ''

	source := flag.String("source", "https://www4.icao.int/doc8643/External/AircraftTypes", "The remote URL where ICAO data can be found.")

	target := flag.String("target", "data/icao.json", "The path to write ICAO aircraft data.")
	stdout := flag.Bool("stdout", false, "Emit ICAO aircraft data to SDOUT.")

	flag.Parse()

	writers := make([]io.Writer, 0)

	fh, err := os.OpenFile(*target, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", *target, err)
	}

	writers = append(writers, fh)

	if *stdout {
		writers = append(writers, os.Stdout)
	}

	wr := io.MultiWriter(writers...)

	data := strings.NewReader("")

	req, err := http.NewRequest("POST", *source, data)

	if err != nil {
		log.Fatalf("Failed to create new request for '%s', %v", *source, err)
	}

	req.Header.Set("Connection", "keep-alive")

	cl := http.Client{}
	rsp, err := cl.Do(req)

	if err != nil {
		log.Fatalf("Failed to request data, %v", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		log.Fatalf("Remote server returned an error, %v", rsp.Status)
	}

	_, err = io.Copy(wr, rsp.Body)

	if err != nil {
		log.Fatalf("Failed to write data, %v", err)
	}
}
