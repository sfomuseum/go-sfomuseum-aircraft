package sfomuseum

import (
	"fmt"
)

type Aircraft struct {
	WOFID          int64  `json:"wof:id"`
	Name           string `json:"wof:name"`
	SFOMuseumID    int    `json:"sfomuseum:aircraft_id"`
	ICAODesignator string `json:"icao:designator,omitempty"`
	WikidataID     string `json:"wd:id,omitempty"`
}

func (a *Aircraft) String() string {
	return fmt.Sprintf("%d %s \"%s\"", a.WOFID, a.ICAODesignator, a.Name)
}
