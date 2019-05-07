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
	$(GO) build -v -o $(BUILD_FOLDER)/src/cmd/filemanager/$(FILE_NAME) ./src/cmd/filemanager/...

.PHONY: src.test
src.test:
	$(GO) test -v ./src/...

.PHONY: src.test-coverage
src.test-coverage:
	$(MKDIR_P) $(BUILD_FOLDER)/src/
	$(GO) test -v -coverprofile=$(BUILD_FOLDER)/src/coverage.txt -covermode=atomic ./src/...
	$(GO) tool cover -html=$(BUILD_FOLDER)/src/coverage.txt -o $(BUILD_FOLDER)/src/coverage.html

## launch apps #############################

define generate_gcr-registry-key
  kubectl create secret docker-registry gcr-registry-key \
    --docker-server=https://gcr.io \
    --docker-username=_json_key \
    --docker-email=chiahsun.jkopay@gmail.com \
    --docker-password='$(shell cat < secret/gcr/$(1).json)' \
    --dry-run -n filemanager -o yaml > deployment/helm/services/charts/secret/templates/gcr.yaml
endef

define generate_filemanager-config
	kubectl create secret generic filemanager-config --from-file=config/$(1).json -n filemanager \
		--dry-run -o yaml > deployment/helm/services/charts/secret/templates/gcs.yaml
endef

.PHONY: apps.init-helm
apps.init-helm:
	helm init
	kubectl create serviceaccount --namespace kube-system tiller
	kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
	kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

.PHONY: apps.install-local
apps.install-local:
	$(call generate_filemanager-config,local)
	yq -y .services deployment/helm/config/local.yaml | helm upgrade --install filemanager-services-local --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/local.yaml | helm upgrade --install filemanager-apps-local --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.install-dev
apps.install-dev:
	$(call generate_gcr-registry-key,develop)
	$(call generate_filemanager-config,develop)
	yq -y .services deployment/helm/config/development.yaml | helm upgrade --install filemanager-services-dev --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/development.yaml | helm upgrade --install filemanager-apps-dev --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.install-stage
apps.install-stage:
	$(call generate_gcr-registry-key,staging)
	$(call generate_filemanager-config,staging)
	yq -y .services deployment/helm/config/staging.yaml | helm upgrade --install filemanager-services-stage --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/staging.yaml | helm upgrade --install filemanager-apps-stage --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.install-rc
apps.install-rc:
	$(call generate_gcr-registry-key,rc)
	$(call generate_filemanager-config,rc)
	yq -y .services deployment/helm/config/rc.yaml | helm upgrade --install filemanager-services-rc --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/rc.yaml | helm upgrade --install filemanager-apps-rc --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.install-prod
apps.install-prod:
	$(call generate_gcr-registry-key,production)
	$(call generate_filemanager-config,production)
	yq -y .services deployment/helm/config/production.yaml | helm upgrade --install filemanager-services-prod --debug -f - deployment/helm/services
	yq -y .apps deployment/helm/config/production.yaml | helm upgrade --install filemanager-apps-prod --debug -f - --set filemanager-server.controller.apiserverImageTag=$(SERVER_VERSION) deployment/helm/apps

.PHONY: apps.teardown-local
apps.teardown-local:
	helm delete --purge filemanager-services-local
	helm delete --purge filemanager-apps-local

.PHONY: apps.teardown-dev
apps.teardown-dev:
	helm delete --purge filemanager-services-dev
	helm delete --purge filemanager-apps-dev

.PHONY: apps.teardown-stage
apps.teardown-stage:
	helm delete --purge filemanager-services-stage
	helm delete --purge filemanager-apps-stage

.PHONY: apps.teardown-rc
apps.teardown-rc:
	helm delete --purge filemanager-services-rc
	helm delete --purge filemanager-apps-rc

.PHONY: apps.teardown-prod
apps.teardown-prod:
	helm delete --purge filemanager-services-prod
	helm delete --purge filemanager-apps-prod

## dockerfiles/ ########################################

.PHONY: dockerfiles.build-local
dockerfiles.build-local:
	docker build --build-arg CONFIG=config/local.json --tag yunchen/file-manager:$(SERVER_VERSION) .

.PHONY: dockerfiles.build-dev
dockerfiles.build-dev:
	docker build --build-arg CONFIG=config/develop.json --tag gcr.io/jello-test-222701/file-manager:$(SERVER_VERSION) .
	docker push gcr.io/jello-test-222701/file-manager:$(SERVER_VERSION)

.PHONY: dockerfiles.build-stage
dockerfiles.build-stage:
	docker build --build-arg CONFIG=config/staging.json --tag gcr.io/jello-stage-223210/file-manager:$(SERVER_VERSION) .
	docker push gcr.io/jello-stage-223210/file-manager:$(SERVER_VERSION)

.PHONY: dockerfiles.build-rc
dockerfiles.build-rc:
	docker build --build-arg CONFIG=config/rc.json --tag gcr.io/jello-stage-223210/file-manager:$(SERVER_VERSION) .
	docker push gcr.io/jello-stage-223210/file-manager:$(SERVER_VERSION)

.PHONY: dockerfiles.build-prod
dockerfiles.build-prod:
	docker build --build-arg CONFIG=config/production.json --tag gcr.io/jello-000001/file-manager:$(SERVER_VERSION) .
	docker push gcr.io/jello-000001/file-manager:$(SERVER_VERSION)

## git tag version ########################################

.PHONY: push.tag
push.tag:
	@echo "Current git tag version:"$(SERVER_VERSION)
	git tag $(SERVER_VERSION)
	git push --tags
