# github-notifier
GitHub-Notifier is a simple personal GitHub timeline notifier, which upon clicking takes you to GitHub page behind the notifcation.

## Usage
```
./github-notifier -u '<username>' -p '<password>'
```

You can run it in a background. It polls every 30 seconds and makes sure the notification are not too old and does not allow repeats.

## Screenshots
Notification Support
![Simple Notification](https://dl.dropboxusercontent.com/u/7162902/github-notifier.png)

Notification Center Support
![Notification Event](https://dl.dropboxusercontent.com/u/7162902/github-notifier-2.png)

## Todo
- [ ] Rewrite terminal-notifier in swift and wrap for go
- [ ] Expose date property in terminal-notifier
- [ ] Support more Event types from GitHub
- [ ] Display partial text of comments and mentions
- [ ] Add Linux notification support

### Note:
This is a heavily changed fork of [deckarep/gosx-notifier](https://github.com/deckarep/gosx-notifier)
