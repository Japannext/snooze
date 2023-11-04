# Snooze v2

# Design choices

* Backend in golang.
* One DB call to write a log upon receiving it.
* Backend in ScyllaDB, denormalized.
* Limited support of search in the web UI (scylladb `LIKE` operator)
