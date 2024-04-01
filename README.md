# Snooze v2

# Design choices

* Backend in golang.
* One DB call to write a log upon receiving it.
* Backend in ScyllaDB, denormalized.
* Limited support of search in the web UI (scylladb `LIKE` operator)

# Components

* Frontend (`snooze-apiserver`):
  - `apiserver/`: Serve the API, expose it to the web interface. Mainly for human user/automation.
  - `ui/`: The Javascript frontend, service by `snooze-apiserver`
* Alert sources:
  - `sources/otel/`: Opentelemetry input. Accept opentelemetry logs as input alerts.
  - `sources/alertmanager/`: Function like a drop-in replacement to Prometheus AlertManager. Same API, accept
    Prometheus alerts.
* Notifiers:
  - `notifier/googlechat/`: Notify alerts to googlechat.
  - `notifier/browser/`: Notify alerts to the browser's open sessions.
  - `notifier/mail/`: Notify alerts by mail
* Shared:
  - `common/`: shared libraries / API
