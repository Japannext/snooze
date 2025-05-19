# snooze-syslog

The `snooze-syslog` component listen on UDP and TCP syslog, parse the received logs, and queue them to the process queue.

This component is organized into the following goroutines:
* `server`: Receive the logs into the receiving queue (channel).
* `parser`: Parse the logs into snooze format, and add them to the publishing queue (channel).
* `publisher`: Batch the logs and publish them to the process queue (rabbitmq).
