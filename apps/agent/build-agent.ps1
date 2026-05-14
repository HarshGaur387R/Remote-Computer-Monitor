# Build 64-bit Agent Binary
# This script builds the rcma.go agent for 64-bit Windows architecture

# Set environment variables for 64-bit Windows
$env:GOOS = "windows"
$env:GOARCH = "amd64"

# Get the script directory (agent directory)
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "Building agent for 64-bit Windows (amd64)..." -ForegroundColor Cyan
Write-Host "Output directory: $scriptDir" -ForegroundColor Cyan

# Build the binary
go build -o "$scriptDir\rcma.exe" "$scriptDir\main.go"

# Check if build was successful
if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Build completed!" -ForegroundColor Green
    
    # Display file info
    if (Test-Path "$scriptDir\rcma.exe") {
        $fileInfo = Get-Item "$scriptDir\rcma.exe"
        Write-Host "Binary created: $($fileInfo.FullName)" -ForegroundColor Green
        Write-Host "File size: $([math]::Round($fileInfo.Length / 1MB, 2)) MB" -ForegroundColor Green
        Write-Host "[SUCCESS] Agent binary built for 64-bit Windows (amd64)" -ForegroundColor Green
    }
} else {
    Write-Host "[ERROR] Build failed with exit code: $LASTEXITCODE" -ForegroundColor Red
    exit $LASTEXITCODE
}