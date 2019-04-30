package sfomuseum

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "log"
	"strconv"
	"strings"
	"sync"
)

type Aircraft struct {
	WOFID          int64  `json:"wof:id"`
	Name           string `json""wof:name"`
	SFOMuseumID    int    `json:"sfomuseum:aircraft_id"`
	ICAODesignator string `json:icao:designator,omitempty"`
	WikidataID     string `json:"wd:id,omitempty"`
}

var lookup_table *sync.Map
var lookup_init sync.Once

type Lookup struct {
}

func NewLookup() (*Lookup, error) {

	var lookup_err error

	lookup_func := func() {

		var aircraft []*Aircraft

		err := json.Unmarshal([]byte(AircraftData), &aircraft)

		if err != nil {
			lookup_err = err
			return
		}

		table := new(sync.Map)

		for idx, craft := range aircraft {

			pointer := fmt.Sprintf("pointer:%d", idx)
			table.Store(pointer, craft)

			str_wofid := strconv.FormatInt(craft.WOFID, 10)

			possible_codes := []string{
				craft.ICAODesignator,
				str_wofid,
			}

			for _, code := range possible_codes {

				if code == "" {
					continue
				}

				pointers := make([]string, 0)
				has_pointer := false

				others, ok := table.Load(code)

				if ok {

					pointers = others.([]string)
				}

				for _, dupe := range pointers {

					if dupe == pointer {
						has_pointer = true
						break
					}
				}

				if has_pointer {
					continue
				}

				pointers = append(pointers, pointer)
				table.Store(code, pointers)
			}

			idx += 1
		}

		lookup_table = table
	}

	lookup_init.Do(lookup_func)

	if lookup_err != nil {
		return nil, lookup_err
	}

	l := Lookup{}
	return &l, nil
}

func (l *Lookup) Find(code string) ([]*Aircraft, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, errors.New("Not found")
	}

	aircraft := make([]*Aircraft, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, errors.New("Invalid pointer")
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, errors.New("Invalid pointer")
		}

		aircraft = append(aircraft, row.(*Aircraft))
	}

	return aircraft, nil
}
