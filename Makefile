## File name for executing
FILE_NAME=file_manager
## Folder content generated files
BUILD_FOLDER = ./build
## command
GO           = go
MKDIR_P      = mkdir -p

################################################

.PHONY: all
all: build-linux test

.PHONY: build-darwin
build-darwin:
	go build -o $(FILE_NAME) src/cmd/gcs-server/main.go

.PHONY: build-linux
build-linux:
	export GOOS=linux && export GOARCH=amd64 &&	go build -o $(FILE_NAME) src/cmd/gcs-server/main.go

.PHONY: test
test: build-linux
	$(MAKE) src.test

.PHONY: run
run:
	go run src/cmd/gcs-server/main.go

.PHONY: zip
zip:
	zip -r config.zip ./config/

## src/ ########################################

.PHONY: src.test
src.test:
	$(GO) test -v ./src/...

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	$(GO) test -v -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html

.PHONY: src.test-coverage-minikube
src.test-coverage-minikube:
	sed -i.bak "s/{{ projectId }}/$(PROJECTID)/g; s/{{ privateKeyId }}/$(PRIVATEKEYID)/g; s#{{ privateKey }}#$(PRIVATEKEY)#g; s/{{ clientEmail }}/$(CLIENTEMAIL)/g; s/{{ clientId }}/$(CLIENTID)/g; s#{{ clientCert }}#$(CLIENTCERT)#g; s/{{ jwtSecretKey }}/$(JWTSECRETKEY)/g;" config/testing.json
	$(MAKE) src.test-coverage
	mv config/testing.json.bak config/testing.json
