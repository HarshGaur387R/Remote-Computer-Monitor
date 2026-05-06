# Install or update rsrc tool for embedding Windows resources
Write-Host "Setting up rsrc tool..." -ForegroundColor Cyan
go install github.com/akavel/rsrc@latest

# Generate rsrc.syso from manifest file
Write-Host "Embedding manifest file..." -ForegroundColor Cyan
& "$(go env GOPATH)\bin\rsrc.exe" -manifest rcmai.exe.manifest -o rsrc.syso

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to embed manifest" -ForegroundColor Red
    exit 1
}

# Build the GUI installer
Write-Host "Building GUI installer..." -ForegroundColor Cyan
go build -ldflags="-H windowsgui" -o ./build/guiinstaller.exe .

# Check if build was successful before running
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful! Running installer..." -ForegroundColor Green
    # Change to build directory and run the executable
    Set-Location build
    .\guiinstaller.exe
} else {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}