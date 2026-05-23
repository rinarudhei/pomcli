<img width="1887" height="978" alt="image" src="https://github.com/user-attachments/assets/354693f8-1cac-4a03-a932-f7e39c309170" />

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

sample git hook script

```
#!/bin/bash
# Pomodoro CLI - Git Post-Commit Hook
# Automatically logs commits to pomodoro sessions

# Exit if pomcli is not installed
if ! command -v pomcli &>/dev/null; then
  exit 0 # Silently exit if pomcli not available
fi

# Get repository name (basename of the git root directory)
REPO_NAME=$(basename "$(git rev-parse --show-toplevel)")

# Get current branch name
BRANCH_NAME=$(git branch --show-current)

# Get the last commit message (title - first line)
COMMIT_TITLE=$(git log -1 --pretty=%s)

# Escape special characters in commit title for command line
# Remove newlines, truncate to reasonable length, and escape quotes
ESCAPED_TITLE=$(echo "$COMMIT_TITLE" | tr '\n' ' ' | head -c 200 | sed "s/'/'\\\\''/g")

# Call pomcli activity-hook
pomcli activity-hook "$REPO_NAME\\$BRANCH_NAME =>" "$ESCAPED_TITLE"

# Exit with pomcli's exit code
exit $?

```

make the script as executable and call it in git config file. (~/.gitconfig)
```

[core]
	hooksPath = /home/{USERNAME}/.config/.git/hooks/

```



