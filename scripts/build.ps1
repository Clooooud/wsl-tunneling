param(
    [string]$Output = "bin\wsl-tunneling.exe",
    [switch]$GUI,
    [switch]$Console
)

$ErrorActionPreference = "Stop"

if ($GUI -and $Console) {
    throw "Use either -GUI or -Console, not both."
}

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

function Invoke-GoBuild {
    param(
        [string]$BuildOutput,
        [string]$Ldflags
    )

    New-Item -ItemType Directory -Force -Path (Split-Path -Parent $BuildOutput) | Out-Null
    if ($Ldflags) {
        & $goPath build -trimpath -ldflags $Ldflags -o $BuildOutput .\cmd\wsl-tunneling
    } else {
        & $goPath build -trimpath -o $BuildOutput .\cmd\wsl-tunneling
    }
}

if ($GUI) {
    Invoke-GoBuild -BuildOutput $Output -Ldflags "-H windowsgui"
} else {
    Invoke-GoBuild -BuildOutput $Output -Ldflags ""
}
