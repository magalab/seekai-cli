$ErrorActionPreference = "Stop"

$Repo = if ($env:SEEKAI_REPO) { $env:SEEKAI_REPO } else { "magalab/seekai-cli" }
$Version = if ($env:SEEKAI_VERSION) { $env:SEEKAI_VERSION } else { "latest" }
$InstallDir = if ($env:SEEKAI_INSTALL_DIR) { $env:SEEKAI_INSTALL_DIR } else { Join-Path $env:USERPROFILE "bin" }
$Asset = "seekai_windows_amd64.exe"

if ($Version -eq "latest") {
  $Url = "https://github.com/$Repo/releases/latest/download/$Asset"
} else {
  $Url = "https://github.com/$Repo/releases/download/$Version/$Asset"
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$Target = Join-Path $InstallDir "seekai.exe"

Write-Host "downloading $Url"
Invoke-WebRequest -Uri $Url -OutFile $Target

Write-Host "installed $Target"
Write-Host "Add $InstallDir to PATH if it is not already present."
Write-Host "run: seekai --help"
