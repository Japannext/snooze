# Requirements

You will need a kubernetes cluster with:
* a working ingress class
* a working storage class
* a configured cert-manager
* opensearch-operator
* dragonflydb-operator (for a redis-like database)

# Getting started

1) Clone the git:
```bash
git clone https://github.com/japannext/snooze
cd snooze/
```

2) Install [taskfile](https://taskfile.dev/installation/)

## Backend development

1) Get [gvm](https://github.com/moovweb/gvm)
2) Install the go version required
```console
> grep 'go 1' go.mod
go 1.23

> gvm install go1.23
[...]
```
3) Build snooze (locally)
```bash
task build
```

## Web development

You need a running backend and [devspace](https://github.com/devspace-sh/devspace).

Important files used by each route:
* Frontend
  - Page: `ui/src/views/*.vue`
  - Type and HTTP call: `ui/src/api/*.vue`
* Backend
  - Types: `pkg/models/*.vue`
  - HTTP route handler: `pkg/apiserver/routes/*.vue`

> There is a relatively consistent name between each files for each route
> in all these directories.

## Process related

The `snooze-process` is executing several sub-components one by one.
* `pkg/processor/*/process.go`
