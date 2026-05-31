# Build on Windows (native amd64). Run from repo root or any path.
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$Backend = Join-Path $Root "backend"
$Bin = Join-Path $Root "bin"
$Out = Join-Path $Bin "eh_adbe_backend-windows-amd64.exe"

New-Item -ItemType Directory -Force -Path $Bin | Out-Null
Push-Location $Backend
try {
    $env:CGO_ENABLED = "0"
    go build -ldflags="-s -w" -o $Out .
    Write-Host "built -> $Out"
} finally {
    Pop-Location
}
