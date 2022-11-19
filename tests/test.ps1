Set-Location $PSScriptRoot\..\

go build -ldflags="-X 'main.Version=v2.0.0'" -o update_v2.0.0.exe .\cmd\example\
Write-Host "## Test Executing Remote Update"
.\update_v2.0.0.exe 

go build -ldflags="-X 'main.Version=v1.0.0'" -o update_v1.0.0.exe .\cmd\example\
Write-Host "## Test Executing Local App"
.\update_v1.0.0.exe 

Write-Host "## Exuting Update"
.\update_v1.0.0.exe -update -updatefile (Resolve-Path .\update_v2.0.0.exe)