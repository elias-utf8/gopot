# gopot

An SSH honeypot written in Go. Gopot listens for incoming SSH connections and presents attackers with a fake interactive shell, logging every session and command to a SQLite database for later analysis. It is designed to blend in as a real Linux system and spoofing a custom OS version string, serving a realistic MOTD, and responding to common shell commands while silently recording the IP address, SSH client version, and full command history of anyone who connects.

## Build & Run

```bash
go build -o bin/gopot ./cmd/gopot/main.go
./bin/gopot
```

The honeypot listens on port `2222` by default. Edit `gopot.toml` to configure the server and shell.

## License

MIT
