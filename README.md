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
* [Alert sources](./sources)
* [Notifiers](./notifiers)
* `common/`: shared libraries / API
