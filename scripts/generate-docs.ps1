#!/usr/bin/env pwsh
# PowerShell version of generate-docs.sh

$workspaceFolder = if ($env:WORKSPACE_FOLDER) { $env:WORKSPACE_FOLDER } else { (Get-Location).Path }

Write-Host "Generating documentation for the project at $workspaceFolder..."
& gomarkdoc -v "$workspaceFolder" --config "$workspaceFolder/.gomarkdoc.yml" -o "$workspaceFolder/README.md"
if ($LASTEXITCODE -ne 0) { Write-Error "gomarkdoc failed for project (exit $LASTEXITCODE)"; exit $LASTEXITCODE }

Write-Host ""
Write-Host "Generating documentation for the 'grapheme' package..."
& gomarkdoc -v "$workspaceFolder/grapheme" --config "$workspaceFolder/.gomarkdoc.yml" -o "$workspaceFolder/grapheme/README.md"
if ($LASTEXITCODE -ne 0) { Write-Error "gomarkdoc failed for grapheme (exit $LASTEXITCODE)"; exit $LASTEXITCODE }

Write-Host ""
Write-Host "Generating documentation for the 'internal/testutils' package..."
& gomarkdoc -v "$workspaceFolder/internal/testutils" --config "$workspaceFolder/.gomarkdoc.yml" -o "$workspaceFolder/internal/testutils/README.md"
if ($LASTEXITCODE -ne 0) { Write-Error "gomarkdoc failed for internal/testutils (exit $LASTEXITCODE)"; exit $LASTEXITCODE }
