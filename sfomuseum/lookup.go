package sfomuseum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-aircraft"
	"github.com/sfomuseum/go-sfomuseum-aircraft/data"
	"io"
	_ "log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var lookup_table *sync.Map
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type SFOMuseumLookupFunc func(context.Context)

type SFOMuseumLookup struct {
	aircraft.Lookup
}

func init() {
	ctx := context.Background()
	aircraft.RegisterLookup(ctx, "sfomuseum", NewLookup)

	lookup_idx = int64(0)
}

// NewLookup will return an `aircraft.Lookup` instance derived from precompiled (embedded) data in `data/sfomuseum.json`
func NewLookup(ctx context.Context, uri string) (aircraft.Lookup, error) {

	fs := data.FS
	fh, err := fs.Open("sfomuseum.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to load data, %v", err)
	}

	lookup_func := NewLookupFuncWithReader(ctx, fh)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

// NewLookup will return an `SFOMuseumLookupFunc` function instance that, when invoked, will populate an `aircraft.Lookup` instance with data stored in `r`.
// `r` will be closed when the `SFOMuseumLookupFunc` function instance is invoked.
// It is assumed that the data in `r` will be formatted in the same way as the procompiled (embedded) data stored in `data/sfomuseum.json`.
func NewLookupFuncWithReader(ctx context.Context, r io.ReadCloser) SFOMuseumLookupFunc {

	lookup_func := func(ctx context.Context) {

		defer r.Close()

		var aircraft []*Aircraft

		dec := json.NewDecoder(r)
		err := dec.Decode(&aircraft)

		if err != nil {
			lookup_init_err = err
			return
		}

		table := new(sync.Map)

		for _, data := range aircraft {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			appendData(ctx, table, data)
		}

		lookup_table = table
	}

	return lookup_func
}

// NewLookupWithLookupFunc will return an `aircraft.Lookup` instance derived by data compiled using `lookup_func`.
func NewLookupWithLookupFunc(ctx context.Context, lookup_func SFOMuseumLookupFunc) (aircraft.Lookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := SFOMuseumLookup{}
	return &l, nil
}

func NewLookupFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) (aircraft.Lookup, error) {

	aircraft_data, err := CompileAircraftData(ctx, iterator_uri, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile aircraft data, %w", err)
	}

	// necessary until there is a NewLookupFuncWithAircraft method
	enc_data, err := json.Marshal(aircraft_data)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal aircraft data, %w", err)
	}

	r := bytes.NewReader(enc_data)
	rc := io.NopCloser(r)

	lookup_func := NewLookupFuncWithReader(ctx, rc)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

func (l *SFOMuseumLookup) Find(ctx context.Context, code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, fmt.Errorf("Code '%s' not found", code)
	}

	aircraft := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, fmt.Errorf("Invalid pointer '%s'", p)
		}

		aircraft = append(aircraft, row.(*Aircraft))
	}

	return aircraft, nil
}

func (l *SFOMuseumLookup) Append(ctx context.Context, data interface{}) error {
	return appendData(ctx, lookup_table, data.(*Aircraft))
}

func appendData(ctx context.Context, table *sync.Map, data *Aircraft) error {

	idx := atomic.AddInt64(&lookup_idx, 1)

	pointer := fmt.Sprintf("pointer:%d", idx)
	table.Store(pointer, data)

	str_wofid := strconv.FormatInt(data.WOFID, 10)
	str_sfomid := strconv.Itoa(data.SFOMuseumID)

	possible_codes := []string{
		data.ICAODesignator,
		str_wofid,
		str_sfomid,
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

	return nil
}
