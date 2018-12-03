## filemanager server version
SERVER_VERSION = latest
## File name for executing
FILE_NAME=file_manager
## Folder content generated files
BUILD_FOLDER = ./build
## command
GO           = go
MKDIR_P      = mkdir -p

################################################

.PHONY: all
all: build build-linux test

.PHONY: build
build:
	$(MAKE) src.build

.PHONY: build-darwin
build-darwin:
	go build -o $(FILE_NAME) src/cmd/filemanger/main.go

.PHONY: build-linux
build-linux:
	export GOOS=linux && export GOARCH=amd64 &&	go build -o $(FILE_NAME) src/cmd/filemanager/main.go

.PHONY: test
test: build-linux
	$(MAKE) src.test

.PHONY: run
run:
	go run src/cmd/filemanager/main.go

.PHONY: zip
zip:
	zip -r config.zip ./config/

## src/ ########################################

.PHONY: src.build
src.build:
	$(GO) build -v ./src/...
	$(MKDIR_P) $(BUILD_FOLDER)/src/cmd/filemanager/
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/filemanager/filemanager ./src/cmd/filemanager/...

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

## launch apps #############################

.PHONY: apps.init-helm
apps.init-helm:
	helm init
	kubectl create serviceaccount --namespace kube-system tiller
	kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
	kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

.PHONY: apps.launch-dev
apps.launch-dev:
	yq -y .services deployment/helm/config/development.yaml | helm install --name filemanager-services-dev --debug --wait -f - deployment/helm/services
	kubectl create configmap filemanager-config --from-file=config/ -n filemanager
	yq -y .apps deployment/helm/config/development.yaml | helm install --name filemanager-apps-dev --debug --wait -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.launch-prod
apps.launch-prod:
	yq -y .services deployment/helm/config/production.yaml | helm install --name filemanager-services-prod --debug --wait -f - deployment/helm/services
	yq -y .apps deployment/helm/config/production.yaml | helm install --name filemanager-apps-prod --debug --wait -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.upgrade-dev
apps.upgrade-dev:
	yq -y .services deployment/helm/config/development.yaml | helm upgrade filemanager-services-dev --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/development.yaml | helm  upgrade filemanager-apps-dev --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.upgrade-prod
apps.upgrade-prod:
	yq -y .services deployment/helm/config/production.yaml | helm upgrade filemanager-services-prod --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/production.yaml | helm  upgrade filemanager-apps-prod --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.teardown-dev
apps.teardown-dev:
	helm delete --purge filemanager-services-dev
	helm delete --purge filemanager-apps-dev
	kubectl delete configmap -n filemanager

.PHONY: apps.teardown-prod
apps.teardown-prod:
	helm delete --purge filemanager-services-prod
	helm delete --purge filemanager-apps-prod

## dockerfiles/ ########################################

.PHONY: dockerfiles.build
dockerfiles.build:
	docker build --tag yunchen/go-gcs:$(SERVER_VERSION) .

## git tag version ########################################

.PHONY: push.tag
push.tag:
	@echo "Current git tag version:"$(SERVER_VERSION)
	git tag $(SERVER_VERSION)
	git push --tags
