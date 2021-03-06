package icao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-aircraft"
	_ "log"
	"strings"
	"sync"
)

type Aircraft struct {
	ModelFullName       string
	Description         string
	WTC                 string
	Designator          string
	ManufacturerCode    string
	AircraftDescription string
	EngineCount         string
	EngineType          string
}

func (a *Aircraft) String() string {
	return fmt.Sprintf("%s %s \"%s\"", a.ManufacturerCode, a.Designator, a.ModelFullName)
}

var lookup_table *sync.Map
var lookup_init sync.Once

type ICAOLookup struct {
	aircraft.Lookup
}

func NewLookup() (aircraft.Lookup, error) {

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

			possible_codes := []string{
				craft.Designator,
				craft.ManufacturerCode,
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

	l := ICAOLookup{}
	return &l, nil
}

func (l *ICAOLookup) Find(code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, errors.New("Not found")
	}

	aircraft := make([]interface{}, 0)

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
