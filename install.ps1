# mr-agent-ai installer for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/SaidHernandez/mr-agent-ai/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$repo    = "SaidHernandez/mr-agent-ai"
$binary  = "mr-agent-ai.exe"
$installDir = "$env:LOCALAPPDATA\mr-agent-ai"

Write-Host ""
Write-Host "  __  __ ____        _                    _         _    ___"
Write-Host " |  \/  |  _ \      / \   __ _  ___ _ __ | |_      / \  |_ _|"
Write-Host " | |\/| | |_) |    / _ \ / _`\|/ _ \ '_ \| __|    / _ \  | |"
Write-Host " | |  | |  _ <    / ___ \ (_| |  __/ | | | |_    / ___ \ | |"
Write-Host " |_|  |_|_| \_\  /_/   \_\__, |\___|_| |_|\__|  /_/   \_\___|"
Write-Host "                          |___/"
Write-Host "  Multi-Agent Skill Installer"
Write-Host ""

# Detect architecture
$arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "32-bit Windows is not supported."
    exit 1
}

# Get latest release version
Write-Host "[1/3] Fetching latest release..."
$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name

# Build download URL
$fileName = "mr-agent-ai_$($version.TrimStart('v'))_windows_$arch.zip"
$url = "https://github.com/$repo/releases/download/$version/$fileName"

# Download
Write-Host "[2/3] Downloading $fileName..."
$tmpZip = Join-Path $env:TEMP $fileName
Invoke-WebRequest -Uri $url -OutFile $tmpZip -UseBasicParsing

# Install
Write-Host "[3/3] Installing to $installDir..."
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}
Expand-Archive -Path $tmpZip -DestinationPath $installDir -Force
Remove-Item $tmpZip

# Add to PATH if not already there
$userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    [System.Environment]::SetEnvironmentVariable("PATH", "$userPath;$installDir", "User")
    Write-Host ""
    Write-Host "[ok] Added $installDir to PATH."
    Write-Host "     Restart your terminal to use mr-agent-ai."
} else {
    Write-Host "[ok] $installDir already in PATH."
}

Write-Host ""
Write-Host "[ok] mr-agent-ai $version installed successfully."
Write-Host "     Run: mr-agent-ai install"
Write-Host ""
