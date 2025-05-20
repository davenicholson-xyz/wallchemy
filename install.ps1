$ErrorActionPreference = "Stop"

$repo = "davenicholson-xyz/wallchemy" 
$binaryName = "wallchemy"
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

$latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$version = $latestRelease.tag_name

$fileName = "$binaryName-windows-$arch-$version.zip"
$downloadUrl = "https://github.com/$repo/releases/download/$version/$fileName"
$tempPath = "$env:TEMP\$fileName"
$extractPath = "$env:TEMP\$binaryName-extracted"

Write-Host "Downloading $fileName..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempPath

# Extract
Write-Host "Extracting..."
Expand-Archive -Path $tempPath -DestinationPath $extractPath -Force

# Install
$installDir = "C:\Program Files\$binaryName"
$exePath = Join-Path $extractPath "$binaryName.exe"

if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}
Copy-Item $exePath -Destination $installDir -Force

# Optionally add to PATH
if (-not ($env:Path -split ";" | Where-Object { $_ -eq $installDir })) {
    Write-Host "Adding to PATH..."
    [Environment]::SetEnvironmentVariable("Path", $env:Path + ";$installDir", [System.EnvironmentVariableTarget]::Machine)
    Write-Host "You may need to restart your terminal or log out and back in to use '$binaryName'."
}

Remove-Item $extractPath -Recurse -Force

Write-Host "✅ $binaryName installed successfully to $installDir"

