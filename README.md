<h1>
  <img src="assets/gopot-icon-transparent.svg" alt="" width="46" align="left" style="margin-right:12px"/>
  gopot
</h1>

An SSH honeypot written in Go.

`gopot` listens for incoming SSH connections and presents attackers with a fake interactive shell, logging every session and command to a SQLite database for later analysis. It blends in as a real Linux system spoofing a custom OS version string, serving a realistic MOTD, and responding to common shell commands while silently recording the IP address, SSH client version, and full command history of anyone who connects.

## Build & run

Build and generate the host key once:

```bash
go build -o bin/gopot ./cmd/gopot/main.go
ssh-keygen -t ed25519 -f host_key -N ""
```

Then run:

```bash
./bin/gopot
```

Listens on port `2223` by default. Edit `gopot.toml` to configure the server and shell.

## License

Apache 2.0
