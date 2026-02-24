# seed.ps1 - Run with: .\scripts\seed.ps1

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$EnvFile = Join-Path $ScriptDir "..\.env"

# Load .env
if (Test-Path $EnvFile) {
    Get-Content $EnvFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.+)$') {
            [System.Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim())
        }
    }
}

$pgAddr = [System.Environment]::GetEnvironmentVariable("POSTGRES_ADDR")

if (-not $pgAddr) {
    Write-Error "ERROR: POSTGRES_ADDR is not set"
    exit 1
}

$seedFile = Join-Path $ScriptDir "..\migrations\seed.sql"

Write-Host "Seeding database..."
psql $pgAddr -f $seedFile
Write-Host "Done."
