package sfomuseum

import (
	"context"
	"github.com/sfomuseum/go-sfomuseum-aircraft"
	"testing"
)

func TestSFOMuseumLookup(t *testing.T) {

	wofid_tests := map[string]int64{
		"B39M": 1528104577,
		"12":   1159289381,
		"A306": 1159289391,
		"17":   1159289391,
		"B744": 1159289915,
	}

	ctx := context.Background()

	lu, err := aircraft.NewLookup(ctx, "sfomuseum://")

	if err != nil {
		t.Fatalf("Failed to create lookup, %v", err)
	}

	for code, wofid := range wofid_tests {

		results, err := lu.Find(ctx, code)

		if err != nil {
			t.Fatalf("Unable to find '%s', %v", code, err)
		}

		if len(results) != 1 {
			t.Fatalf("Invalid results for '%s'", code)
		}

		a := results[0].(*Aircraft)

		if a.WOFID != wofid {
			t.Fatalf("Invalid match for '%s', expected %d but got %d", code, wofid, a.WOFID)
		}
	}
}
