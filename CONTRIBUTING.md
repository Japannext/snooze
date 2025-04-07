# Requirements

You will need the following operators:
* cert-manager
* opensearch-operator
* dragonflydb-operator

You will need the following programs:
* [taskfile](https://taskfile.dev/installation/) (An alternative to Makefile)
* [mise](https://mise.jdx.dev/installing-mise.html) (version manager for golang and bun)
* [devspace](https://github.com/devspace-sh/devspace) (for running the webui auto-reloading in kubernetes)

# Getting started

1) Clone the git:
```bash
git clone https://github.com/japannext/snooze
cd snooze/
```

Install the golang and bun version used by the project:
```bash
mise install
```

## Backend development

Build snooze (locally)
```bash
task build
```

## Web development

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

# Running unit tests

The unit tests need to run against real databases.
The databases can be configured with the same environment
variables that snooze uses, in `.env.local`.

Minimum required:
```
OPENSEARCH_ADDRESSES=http://127.0.0.1:9200
OPENSEARCH_USERNAME=admin
OPENSEARCH_PASSWORD=admin

REDIS_ADDRESS=https://127.0.0.1:6379

NATS_URL=http://127.0.0.1:4222
```

> Note: Opensearch requires a specific kernel parameter that is not always set correctly by default.
> Run `sudo sysctl -w vm.max_map_count=262144` to fix it locally, or on your Kubernetes nodes.

The following options are possible:
* Run it in Kubernetes: `task setup:k8s` (will use the current context, and the `snooze-testing` namespace)
* Run it in Docker: `task setup:docker`
* Use external databases.
