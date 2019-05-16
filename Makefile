fmt:
	go fmt *.go
	go fmt cmd/*.go
	go fmt icao/*.go
	go fmt sfomuseum/*.go

tools:
	build -o bin/lookup cmd/lookup/main.go
