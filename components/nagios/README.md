# Nagios/Icinga CLI plugin

CLI helper to configure Nagios/Icinga notification plugin to send notifications to snooze.

# Nagios


# Icinga

1) Copy the `snooze-nagios` binary to `/usr/lib64/nagios/plugins/`
2) Configure icinga as such:
```
object NotificationCommand "snooze-nagios host" {
    import "plugin-notification-command"
    command = [PluginDir + "/snooze-nagios", "host"]
    arguments = {
        "--url" = "https://snooze.example.com"
        "--hostname" = host.name
        "--timestamp" = "$host.last_check$"
        "--state" = "$host.state$"
        "--icinga" = {}
        "--labels" = "icinga.host.execution_time=$host.execution_time$"
    }
}
```
