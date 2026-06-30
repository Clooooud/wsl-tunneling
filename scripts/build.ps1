param(
    [string]$Output = "bin\wsl-tunneling.exe",
    [switch]$GUI,
    [switch]$CLI
)

$ErrorActionPreference = "Stop"

if ($GUI -and $CLI) {
    throw "Use either -GUI or -CLI, not both."
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
        [string]$Ldflags,
        [string]$Package = ".\cmd\wsl-tunneling"
    )

    New-Item -ItemType Directory -Force -Path (Split-Path -Parent $BuildOutput) | Out-Null
    if ($Ldflags) {
        & $goPath build -trimpath -ldflags $Ldflags -o $BuildOutput $Package
    } else {
        & $goPath build -trimpath -o $BuildOutput $Package
    }
}

if ($GUI) {
    Invoke-GoBuild -BuildOutput $Output -Ldflags "-H windowsgui"
} elseif ($CLI) {
    if ($Output -eq "bin\wsl-tunneling.exe") {
        $Output = "bin\wsl-tunneling-cli.exe"
    }
    Invoke-GoBuild -BuildOutput $Output -Ldflags "" -Package ".\cmd\wsl-tunneling-cli"
} else {
    $outputDir = Split-Path -Parent $Output
    if (-not $outputDir) {
        $outputDir = "."
    }
    Invoke-GoBuild -BuildOutput (Join-Path $outputDir "wsl-tunneling.exe") -Ldflags "-H windowsgui" -Package ".\cmd\wsl-tunneling"
    Invoke-GoBuild -BuildOutput (Join-Path $outputDir "wsl-tunneling-cli.exe") -Ldflags "" -Package ".\cmd\wsl-tunneling-cli"
}
