package template

import (
	"io"
	gotemplate "text/template"
	"time"
)

// the {{ tick }} stuff is because this https://github.com/golang/go/issues/18221#issuecomment-394255883

const aircraftDataTemplate string = `package {{ .Package }}

// this file was generated by the go-sfomuseum-aircraft-tools package on {{ .LastUpdate }}
{{ $tick := "` + "`" + `" }}
const AircraftData string = {{ $tick }}{{ .Data }}{{ $tick }}
`

type AircraftDataVars struct {
	Package    string
	Data       string
	LastUpdate string
}

func RenderAircraftData(wr io.Writer, vars *AircraftDataVars) error {

	now := time.Now()

	vars.LastUpdate = now.Format(time.RFC3339)

	t := gotemplate.New("data")
	t, err := t.Parse(aircraftDataTemplate)

	if err != nil {
		return err
	}

	return t.Execute(wr, vars)
}
