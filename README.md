# wsl-tunneling

`wsl-tunneling` is a small Windows-side controller for a WSL user-mode network tunnel backed by `containers/gvisor-tap-vsock`.

The first implementation targets one configurable WSL distro. It starts `gvforwarder` inside that distro and connects it to `gvproxy.exe` on Windows through the same `stdio:...listen-stdio=accept` transport used by Podman for WSL user-mode networking.

## Current scope

- CLI commands: `doctor`, `start`, `stop`, `restart`, `status`, `daemon`, `logs`, `install-service`, `uninstall-service`.
- Downloads versioned `gvproxy.exe` and `gvforwarder` release assets, or accepts explicit asset URLs.
- Saves WSL DNS and default route state before start.
- Replaces `/etc/resolv.conf` with `nameserver 192.168.127.1` while running.
- Restores DNS and route state on stop.
- Runs a foreground daemon that supervises and restarts the forwarder.

## Requirements

- Windows with WSL installed and updated.
- A target WSL distro with `bash`, `ip`, `awk`, `grep`, `cp`, `kill`, and Windows interop enabled.
- Go 1.22 or newer to build from source.
- Network access to GitHub releases unless `gvproxyUrl` and `gvforwarderUrl` are configured explicitly.

## Build

```powershell
go build -o bin/wsl-tunneling.exe ./cmd/wsl-tunneling
```

## Release

GitHub Actions builds Windows binaries and publishes them to a GitHub Release when a tag starting with `v` is pushed:

```powershell
git tag v0.1.0
git push origin v0.1.0
```

The release contains:

- `wsl-tunneling-windows-amd64.exe`
- `wsl-tunneling-windows-arm64.exe`
- SHA-256 checksum files

## Basic use

Create a config:

```powershell
bin\wsl-tunneling.exe init-config --config $env:APPDATA\wsl-tunneling\config.json
```

Edit the `distro` value, then run diagnostics:

```powershell
bin\wsl-tunneling.exe doctor --config $env:APPDATA\wsl-tunneling\config.json
```

Start the tunnel:

```powershell
bin\wsl-tunneling.exe start --config $env:APPDATA\wsl-tunneling\config.json
```

Check from Windows:

```powershell
bin\wsl-tunneling.exe status --config $env:APPDATA\wsl-tunneling\config.json
```

Check from WSL:

```bash
cat /etc/resolv.conf
ip route
curl https://example.com
```

Stop and restore networking:

```powershell
bin\wsl-tunneling.exe stop --config $env:APPDATA\wsl-tunneling\config.json
```

## Background mode

Install as a Windows scheduled task that starts the daemon at logon:

```powershell
bin\wsl-tunneling.exe install-service --config $env:APPDATA\wsl-tunneling\config.json
```

Remove the scheduled task:

```powershell
bin\wsl-tunneling.exe uninstall-service --config $env:APPDATA\wsl-tunneling\config.json
```

## Recovery

If WSL networking is left in a bad state, first try:

```powershell
bin\wsl-tunneling.exe stop --config $env:APPDATA\wsl-tunneling\config.json
```

If that does not restore the route or DNS, use:

```powershell
wsl --shutdown
```

Then start the target distro normally. The tool stores transient state under `/mnt/wsl/wsl-tunneling` by default.

## Known limits

- This is not yet multi-distro aware.
- ICMP forwarding is limited by `gvisor-tap-vsock`; use TCP checks such as `curl` for validation.
- Port forwarding and `win-sshproxy` are not part of this MVP.
- `install-service` currently installs a Windows scheduled task, not a native Service Control Manager service.
