#!/usr/bin/env pwsh
# PowerShell version of build-examples.sh

function Convert-Size {
    param([long]$bytes)
    if ($bytes -ge 1GB) { "{0:N2} GB" -f ($bytes/1GB) }
    elseif ($bytes -ge 1MB) { "{0:N2} MB" -f ($bytes/1MB) }
    elseif ($bytes -ge 1KB) { "{0:N2} KB" -f ($bytes/1KB) }
    else { "$bytes B" }
}

$workspaceFolder = if ($env:WORKSPACE_FOLDER) { $env:WORKSPACE_FOLDER } else { (Get-Location).Path }
$examplesPath = Join-Path $workspaceFolder 'examples'

Write-Host "Size before build:"
if (-not (Test-Path $examplesPath)) {
    Write-Host "Directory not found: $examplesPath" -ForegroundColor Red
    exit 1
}

Get-ChildItem -Path $examplesPath -Force | Where-Object { $_.Name -eq 'examples' } | ForEach-Object {
    if ($_.PSIsContainer) {
        Write-Host ('{0,-10} {1,12} {2}' -f $_.Mode, '<DIR>', $_.Name)
    } else {
        Write-Host ('{0,-10} {1,12} {2}' -f $_.Mode, $_.Length, $_.Name)
        Write-Host ('{0,-10} {1,12} {2}' -f $_.Mode, $(Convert-Size $_.Length), $_.Name)
    }
}

Write-Host ""
Write-Host "Size after build:"

Push-Location $examplesPath
try {
    & go build --ldflags "-s -w" -o examples
    if ($LASTEXITCODE -ne 0) {
        Write-Host "go build returned exit code $LASTEXITCODE" -ForegroundColor Red
        Pop-Location
        exit $LASTEXITCODE
    }
} catch {
    Write-Host "go build failed: $_" -ForegroundColor Red
    Pop-Location
    exit 1
}

Get-ChildItem -Path . -Force | Where-Object { $_.Name -eq 'examples' } | ForEach-Object {
    Write-Host ('{0,-10} {1,12} {2}' -f $_.Mode, $_.Length, $_.Name)
    Write-Host ('{0,-10} {1,12} {2}' -f $_.Mode, (Convert-Size $_.Length), $_.Name)
}

Pop-Location
