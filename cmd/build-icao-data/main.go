package main

import (
	"flag"
	"github.com/sfomuseum/go-sfomuseum-aircraft-tools/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	// curl 'https://www4.icao.int/doc8643/External/AircraftTypes' -H 'Connection: keep-alive' --data ''

	source := flag.String("source", "https://www4.icao.int/doc8643/External/AircraftTypes", "...")

	flag.Parse()

	data := strings.NewReader("")

	req, err := http.NewRequest("POST", *source, data)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Connection", "keep-alive")

	cl := http.Client{}
	rsp, err := cl.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		log.Fatal(rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		log.Fatal(err)
	}

	vars := template.AircraftDataVars{
		Package: "icao",
		Data:    string(body),
	}

	err = template.RenderAircraftData(os.Stdout, &vars)

	if err != nil {
		log.Fatal(err)
	}
}
