# Snooze roadmap

# V2 roadmap

1) Monorepo
* Components to monorepo:
  * all server components (rule/aggregate/snooze/notify)
  * all input components (syslog/snmptrap/etc)
  * all output components (googlechat/mail/etc)
  * snooze client
* We should have multi-components unit/integration tests.

2) Separate the server into different components
* Meant to be ran and scaled in a kubernetes infrastructure
* General idea: separate process (rule/aggregate/snooze) from notify (notification/action).
  This makes sense, since process is done quickly, while notify might delay the action in a queue, and may even use a different
  data structure. Also, process might want to store records in a fast-to-write database (cassandra, elasticsearch, ...), while the
  flexibility of mongodb might be good enough for notifications.
* router:
  -> Should provide an interface to route to /process and /notify. Meant for the web interface and human interactions.
  Most plugin scripts will talk to /process directly instead of going through the /router. The main goal is to provide a unified interface,
  and do the authentication/authorization.
  -> Should extend the APIs it talks to. Let's make it talk with input/output plugins (googlechat, syslog, etc) so we get better config mgmt/monitoring info.
  -> Should provide authentication/authorization
* process:
  -> Need a db for alerts (mongodb, cassandra, etc)
* notify:
  -> Need a scheduler queue (a consumer which can react to datetime), which can be implemented with different dbs (let's re-use our mongodb implementation)

3) Migrate away from CoreUI?
* There are some good Vue3 components libraries nowadays, with almost the same functionalities as some custom things we wrote
* Let's rewrite things in composition API
* Let's rewrite things in typescript

4) Kubernetes support
* We should have a helm recipe to deploy snooze to kubernetes
* We should integrate with most common tools used in kubernetes (vector, fluentd, etc)
* We should make the alert fields configurable, and provide good defaults for kubernetes.
  E.g. timestamp/namespace/pod/message
* We should provide example of how to send alerts to snooze (using sidecar containers with custom rules and a `severity`, etc)

5) OIDC support
* We should support OpenID connect Authentication method
* We should test it in a kubernetes setup with keycloak (very standard OSS setup)
* We should support authorization just like with ldap (populating the `groups` field?)

6) Moving to yarn
* Yarn is faster
* Yarn is better for CI, because of the way its dependency lock works
* No need to ship the deps in git, we can activate the zero install to get the good lock file, and avoid committing the zero-install
  directory. It gets the advantage of a good lock file, while avoiding to ship dependencies in the git.

7) Moving the Vite
* Vite will make developement faster
* Vite will make packaging faster

8) Recoding part of the code-base in async?
* Async is better than threads, since we get a result when it returns (so we can get rid of the SurvivingThead system),
  and exceptions are handled better as well.
* Different threads interactions are hard to test. Maybe async is easier?

# TODO

* Refactor configuration management (using `python-confuse` maybe)
* Rewrite DB path
* Replace list of Rules with a Root Rules (Ex: To handle global actions such as maintenance)
* Transform a record search into a rule
* Export and import DB
* When restarting snooze server, should not have to restart syslog as well
* Client CLI
  * Client CLI dynamically updated with new versions of the API schema
* Plugin manager (wraps Pip)
* Play with montydb (remove tinydb, replace mongomock)
* Time constraints: holidays (use a custom calendar?)
* Time constraints in rules
* Personal environment
* Recode severity to map fixed values

## Long term

* TCP stream for resources that can be updated. On update, snooze-server should send a
push notification to all listening clients.
* Replace http basic auth with digest auth
* Vagrant testing
* LDAP backend using SASL (GSSAPI)
* Auto-generated certificates
* Refactor testing to use the same samples for all tests
  * Standardize testing for modules
* Use tox for supporting multiple versions of python
* When clicking on a rule row, redirect to Records and apply a search listing all records matching this rule (and all parents conditions as well)
* Dedicated view for each record by clicking on it (or have a button)
  * Basically replace "More" by a dedicated view (remove it or rename it)
* When filling up a Condition, make it more user friendly by showing a dropdown of suggestions as you type
* Snooze notify
  * Be able to optionally assign one or multiple commands when creating a snooze filter (using a dropdown like for notifications)
  * Before 'Abort and Write to DB' when a record is snoozed, run these commands
* Add auto documentation for each known error log received. Possibility for the user to add more
* Put "comment, shelve, ack, close" buttons in a div overlapping messages and being shown when hovering
