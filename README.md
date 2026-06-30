# wsl-tunneling

`wsl-tunneling` is a Windows tray app for running a WSL user-mode network tunnel backed by `gvisor-tap-vsock`.

It is meant for the daily workflow where WSL needs to send traffic through a Windows-side tunnel while keeping DNS and the default route recoverable. You can start and stop the tunnel from the tray icon, enable startup at logon, and use PowerShell commands when you need diagnostics.

## What It Changes

When the tunnel is started, the app:

- starts `gvforwarder` in the configured WSL distro and connects it to `gvproxy.exe` on Windows;
- creates the WSL-side `wsl-tunneling` network interface;
- changes the WSL default route to `192.168.127.1`;
- writes `/etc/resolv.conf` with `nameserver 192.168.127.1` and your Windows DNS search suffixes;
- saves the previous WSL route, DNS, and WSL network config so `Stop` can restore them.

WSL shares its network namespace across distros, so the network interface can be visible from multiple distros. The `distro` setting is the distro where the helper process is launched, not an isolated per-distro network target.

## Requirements

- Windows with WSL installed and updated.
- A WSL distro with `bash`, `ip`, `awk`, `grep`, `cp`, `kill`, and Windows interop enabled.
- Network access to GitHub releases unless you configure explicit `gvproxyUrl` and `gvforwarderUrl` values.

## Install

Download the Windows zip from the GitHub Release for your architecture:

- `wsl-tunneling-windows-amd64.zip`
- `wsl-tunneling-windows-arm64.zip`

Extract the zip somewhere stable. It contains two files:

- `wsl-tunneling.exe`: the tray app. Double-click this one.
- `wsl-tunneling-cli.exe`: the PowerShell/cmd command app.

For example:

```powershell
mkdir $env:LOCALAPPDATA\wsl-tunneling
Expand-Archive .\wsl-tunneling-windows-amd64.zip $env:LOCALAPPDATA\wsl-tunneling -Force
```

Double-click `wsl-tunneling.exe` to open the tray app. No terminal window should stay open.

## First Setup

Open PowerShell in the folder containing the two executables, then create the default config:

```powershell
.\wsl-tunneling-cli.exe init-config
```

The config is written to:

```text
%APPDATA%\wsl-tunneling\config.json
```

Open that file and set `distro` to the WSL distro where the helper should run. For example:

```json
{
	"distro": "Ubuntu-24.04"
}
```

The generated config includes the other defaults, so usually only `distro` needs attention.

Run a quick check:

```powershell
.\wsl-tunneling-cli.exe doctor
```

## Use The Tray App

Run the executable without arguments, or double-click it:

```powershell
.\wsl-tunneling.exe
```

Right-click the tray icon to open the menu:

- `Start` starts the WSL tunnel.
- `Stop` stops it and restores WSL route/DNS state.
- `Open config folder` opens `%APPDATA%\wsl-tunneling` and creates an example config if missing.
- `Start on boot` toggles launching the tray app at Windows logon.
- `Quit` closes the tray process.

To enable startup at logon from PowerShell:

```powershell
.\wsl-tunneling-cli.exe install-service
```

This registers `wsl-tunneling.exe`, the tray app, at Windows logon. Keep both executables in the same folder.

To remove startup at logon:

```powershell
.\wsl-tunneling-cli.exe uninstall-service
```

## Useful Commands

Start the tunnel:

```powershell
.\wsl-tunneling-cli.exe start
```

Check status:

```powershell
.\wsl-tunneling-cli.exe status
```

Stop and restore WSL networking:

```powershell
.\wsl-tunneling-cli.exe stop
```

Restart:

```powershell
.\wsl-tunneling-cli.exe restart
```

Show logs:

```powershell
.\wsl-tunneling-cli.exe logs
```

You can also override the config path when testing:

```powershell
.\wsl-tunneling-cli.exe status --config C:\path\to\config.json
```

## Check From WSL

After starting the tunnel, useful checks inside WSL are:

```bash
ip route show default
cat /etc/resolv.conf
curl https://example.com
```

You should see a default route through `192.168.127.1` on the `wsl-tunneling` interface, and `/etc/resolv.conf` should contain `nameserver 192.168.127.1` plus DNS search suffixes copied from Windows.

## Configuration

The default config path is `%APPDATA%\wsl-tunneling\config.json`.

Common fields:

- `distro`: WSL distro used to run the helper process.
- `gvisorVersion`: `gvisor-tap-vsock` release version to download.
- `interfaceName`: WSL network interface name, default `wsl-tunneling`.
- `gatewayIp`: gateway exposed by the tunnel, default `192.168.127.1`.
- `deviceIp`: WSL-side address, default `192.168.127.2`.
- `dnsSearchSuffixes`: optional explicit DNS search suffixes. Leave empty to use Windows DNS suffixes automatically.
- `terminateOnStop`: whether `stop` should terminate the configured WSL distro after cleanup.

Most users only need to set `distro`.

## Recovery

If WSL networking is left in a bad state, first run:

```powershell
.\wsl-tunneling-cli.exe stop
```

If DNS or routes are still wrong, restart WSL:

```powershell
wsl --shutdown
```

Then start your distro again normally.

## Known Limits

- WSL networking is shared across distros, so the tunnel is not isolated per distro.
- ICMP forwarding is limited by `gvisor-tap-vsock`; use TCP checks such as `curl` for validation.
- Port forwarding and `win-sshproxy` are not part of this MVP.
- `install-service` installs a Windows logon startup entry, not a native Service Control Manager service.

## Build From Source

You only need this section if you are developing or building the executable yourself.

Requirements:

- Go 1.22 or newer.

Build both executables:

```powershell
.\scripts\build.ps1
```

This creates:

- `bin\wsl-tunneling.exe`: GUI tray app.
- `bin\wsl-tunneling-cli.exe`: console command app.

To build only one side:

```powershell
.\scripts\build.ps1 -GUI -Output bin\wsl-tunneling.exe
.\scripts\build.ps1 -CLI -Output bin\wsl-tunneling-cli.exe
```

Run tests:

```powershell
go test ./...
```

## Release Notes For Maintainers

GitHub Actions builds Windows binaries and publishes them to a GitHub Release when a tag starting with `v` is pushed:

```powershell
git tag v0.1.0
git push origin v0.1.0
```

The release contains one zip per Windows architecture. Each zip contains `wsl-tunneling.exe` and `wsl-tunneling-cli.exe`, plus SHA-256 checksum files for the zips.
