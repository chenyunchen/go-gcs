EXE=file_manager

all: build-linux

build-darwin:
	go build -o $(EXE) src/cmd/gcs-server/main.go

build-linux:
	export GOOS=linux && export GOARCH=amd64 &&	go build -o $(EXE) src/cmd/gcs-server/main.go

run:
	go run src/cmd/gcs-server/main.go
