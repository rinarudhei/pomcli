# pomcli

Opinionated Pomodoro counter CLI application. Focus booster.

## Pomodoro App

Usage:
pomcli [flags]

Flags:
-h, --help help for pomcli
-l, --long-break duration Long break duration (default 59m0s)
-p, --pomodoro duration Pomodoro duration (default 50m0s)
-s, --short-break duration Short break duration (default 10m0s)
-v, --version version for pomcli

## Insert Commit or Activity Log

Utilize git hook to update activity in Pomcli asynchronously
or manually using this command:

```
pomcli activity-hook ...[args]

//example: pomcli activity-hook {project-name}\{branch-name} {description}
```
