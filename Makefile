NAME := snooze
REPO := ghcr.io/japannext/snooze

VERSION := $(shell grep -oP 'appVersion: \K.*' charts/snooze/Chart.yaml)
CHART_VERSION := $(shell grep -oP 'appVersion: \K.*' charts/snooze/Chart.yaml)
COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/v2/log.proto

.PHONY: process-build
process-build:
	go build cmd/process/main.go

.PHONY: process-develop
process-develop:
	docker build ./process.Dockerfile -t ${LOCAL_REPO}/snooze-process:develop --build-arg VERSION=$(VERSION)-develop --build-arg COMMIT=$(COMMIT)
	docker push ${LOCAL_REPO}/snooze-process:develop

.PHONY: chart-develop
chart-develop:
	helm package charts/$(NAME) --app-version develop --version 0.0.0-dev -d .charts/
	helm cm-push .charts/$(NAME)-0.0.0-dev.tgz jnx-repo-upload

.PHONY: helmfile-develop
helmfile-develop:
	helmfile -f .helmfile.yaml sync

.PHONY: develop
develop: process-develop chart-develop helmfile-develop
