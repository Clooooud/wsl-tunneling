# wsl-tunneling

`wsl-tunneling` is a small Windows-side controller for a WSL user-mode network tunnel backed by `containers/gvisor-tap-vsock`.

The first implementation targets one configurable WSL distro. It starts `gvforwarder` inside that distro and connects it to `gvproxy.exe` on Windows through the same `stdio:...listen-stdio=accept` transport used by Podman for WSL user-mode networking.

## Current scope

- CLI commands: `doctor`, `start`, `stop`, `restart`, `status`, `tray`, `daemon`, `logs`, `install-service`, `uninstall-service`.
- Downloads versioned `gvproxy.exe` and `gvforwarder` release assets, or accepts explicit asset URLs.
- Saves WSL DNS and default route state before start.
- Replaces `/etc/resolv.conf` with `nameserver 192.168.127.1` while running.
- Restores DNS and route state on stop.
- Runs a Windows tray controller with start, stop, config folder, and start-on-boot actions.
- Keeps a foreground daemon command available for console supervision and troubleshooting.

## Requirements

- Windows with WSL installed and updated.
- A target WSL distro with `bash`, `ip`, `awk`, `grep`, `cp`, `kill`, and Windows interop enabled.
- Go 1.22 or newer to build from source.
- Network access to GitHub releases unless `gvproxyUrl` and `gvforwarderUrl` are configured explicitly.

## Build

```powershell
.\scripts\build.ps1
```

The default build creates `bin\wsl-tunneling.exe`, a single executable that works both ways: double-clicking it opens the tray and hides the private console window Windows creates for Explorer launches, while running it from PowerShell keeps normal console output for commands such as `status`, `doctor`, and `start`.

For a tray-only GUI build, use:

```powershell
.\scripts\build.ps1 -GUI -Output bin\wsl-tunneling.exe
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

The release binaries are console-capable apps that hide their private console window when launched in tray mode from Explorer.

## Basic use

The default config path is `%APPDATA%\wsl-tunneling\config.json`. Use `--config` only when you want a different file.

Create a config:

```powershell
bin\wsl-tunneling.exe init-config
```

Edit the `distro` value, then run diagnostics:

```powershell
bin\wsl-tunneling.exe doctor
```

Start the tunnel:

```powershell
bin\wsl-tunneling.exe start
```

Check from Windows:

```powershell
bin\wsl-tunneling.exe status
```

Check from WSL:

```bash
cat /etc/resolv.conf
ip route
curl https://example.com
```

Stop and restore networking:

```powershell
bin\wsl-tunneling.exe stop
```

## Background mode

Run the binary without a command to open the tray controller:

```powershell
bin\wsl-tunneling.exe
```

`bin\wsl-tunneling.exe tray` does the same thing explicitly from a console.

The tray menu can start and stop the tunnel, open the config folder, toggle `Start on boot`, and quit the tray process. If the config file does not exist yet, `Open config folder` creates an example config before opening the folder.

Install as a Windows scheduled task that starts the tray at logon:

```powershell
bin\wsl-tunneling.exe install-service
```

Remove the scheduled task:

```powershell
bin\wsl-tunneling.exe uninstall-service
```

## Recovery

If WSL networking is left in a bad state, first try:

```powershell
bin\wsl-tunneling.exe stop
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
- `install-service` installs a Windows scheduled task, not a native Service Control Manager service.
