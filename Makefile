cli:
	go build -mod vendor -o bin/build-icao-data cmd/build-icao-data/main.go
	go build -mod vendor -o bin/build-sfomuseum-data cmd/build-sfomuseum-data/main.go
	go build -mod vendor -o bin/lookup cmd/lookup/main.go

rebuild:
	go build -mod vendor -o bin/build-icao-data cmd/build-icao-data/main.go
	go build -mod vendor -o bin/build-sfomuseum-data cmd/build-sfomuseum-data/main.go
	bin/build-icao-data
	bin/build-sfomuseum-data
	go build -mod vendor -o bin/lookup cmd/lookup/main.go
