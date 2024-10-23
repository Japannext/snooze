# Syslog

The syslog component listen on UDP and TCP for syslog message.

# Configuration

## Using rsyslogd

When using a rsyslogd relay, we recommend this configuration to
send error logs to snooze:
```conf
ruleset(name="snooze") {
    action(
        name="snooze1"
        type="omfwd"
        protocol="tcp"
        target="snooze-syslog.example.com"
        port="514"
        queue.type="linkedlist"
        queue.filename="snooze1"
        queue.saveonshutdown="on"
        queue.size="100000"
    )
}

# Send to snooze when the severity is error (or worse)
if $syslogseverity <= 3 then call snooze
```
