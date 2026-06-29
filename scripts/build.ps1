param(
    [string]$Output = "bin\wsl-tunneling.exe",
    [switch]$GUI
)

$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path (Split-Path -Parent $Output) | Out-Null

$go = Get-Command go -ErrorAction SilentlyContinue
if ($null -eq $go -and (Test-Path "$env:USERPROFILE\go\bin\go.exe")) {
    $go = Get-Item "$env:USERPROFILE\go\bin\go.exe"
}

if ($null -eq $go) {
    throw "go.exe was not found in PATH or $env:USERPROFILE\go\bin"
}

if ($go.PSObject.Properties.Name -contains "Source" -and $go.Source) {
    $goPath = $go.Source
} else {
    $goPath = $go.FullName
}

$ldflags = ""
if ($GUI) {
    $ldflags = "-H windowsgui"
}

if ($ldflags) {
    & $goPath build -trimpath -ldflags $ldflags -o $Output .\cmd\wsl-tunneling
} else {
    & $goPath build -trimpath -o $Output .\cmd\wsl-tunneling
}
