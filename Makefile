fmt:
	go fmt *.go
	go fmt cmd/lookup/main.go
	go fmt icao/*.go
	go fmt sfomuseum/*.go

tools:
	go build -o bin/lookup cmd/lookup/main.go
