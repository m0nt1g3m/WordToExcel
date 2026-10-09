param(
    [Parameter(Mandatory = $true)]
    [string]$SourceExePath,

    [Parameter(Mandatory = $true)]
    [string]$OutputInstallerPath,

    [Parameter(Mandatory = $true)]
    [string]$AppName,

    [Parameter(Mandatory = $false)]
    [string]$AppVersion = ""
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($AppVersion)) {
    $AppVersion = "0.0.0-dev"
}

if ($AppVersion.StartsWith("v")) {
    $AppVersion = $AppVersion.Substring(1)
}

if (-not (Test-Path $SourceExePath)) {
    throw "Source EXE not found: $SourceExePath"
}

$buildDir = Split-Path -Parent $OutputInstallerPath
$installerSourceDir = Join-Path $buildDir "installer_source"
$installerSourceExe = Join-Path $installerSourceDir "$AppName.exe"
$issPath = Join-Path $buildDir "$AppName-installer.iss"

New-Item -ItemType Directory -Path $installerSourceDir -Force | Out-Null
Copy-Item $SourceExePath $installerSourceExe -Force

$installerBaseName = [System.IO.Path]::GetFileNameWithoutExtension($OutputInstallerPath)
if ([string]::IsNullOrWhiteSpace($installerBaseName)) {
    throw "Installer output filename must have a non-empty basename: $OutputInstallerPath"
}

$iconPath = Join-Path $PSScriptRoot "..\icons\icon.ico"
if (-not (Test-Path $iconPath)) {
    $iconPath = Join-Path $PSScriptRoot "..\icons\icon_win.png"
}

$issContent = @"
[Setup]
AppName=$AppName
AppVersion=$AppVersion
DefaultDirName={autopf}\$AppName
DefaultGroupName=$AppName
OutputDir=$buildDir
OutputBaseFilename=$installerBaseName
Compression=lzma
SolidCompression=yes
PrivilegesRequired=lowest
SetupIconFile=$iconPath
UninstallDisplayIcon={app}\$AppName.exe
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64

[Files]
Source: "$installerSourceExe"; DestDir: "{app}"

[Icons]
Name: "{group}\$AppName"; Filename: "{app}\$AppName.exe"
"@

Set-Content -Path $issPath -Value $issContent -Encoding UTF8

$innoSetupCommand = Get-Command iscc -ErrorAction SilentlyContinue
if (-not $innoSetupCommand) {
    Write-Host "Inno Setup not found. Installing via winget..." -ForegroundColor Yellow
    winget install --id JRSoftware.InnoSetup -e --accept-source-agreements --accept-package-agreements --disable-interactivity
    $innoSetupCommand = Get-Command iscc -ErrorAction SilentlyContinue
}

if (-not $innoSetupCommand) {
    throw "Inno Setup is not available and could not be installed."
}

& $innoSetupCommand.Source $issPath

$generatedInstaller = Join-Path $buildDir "$installerBaseName.exe"
if (-not (Test-Path $generatedInstaller)) {
    throw "Installer was not created: $generatedInstaller"
}

$generatedInstallerFullPath = [System.IO.Path]::GetFullPath($generatedInstaller)
$outputInstallerFullPath = [System.IO.Path]::GetFullPath($OutputInstallerPath)
if (-not [string]::Equals($generatedInstallerFullPath, $outputInstallerFullPath, [System.StringComparison]::OrdinalIgnoreCase)) {
    Copy-Item $generatedInstallerFullPath $outputInstallerFullPath -Force
}

Write-Host "✅ Windows installer created: $OutputInstallerPath" -ForegroundColor Green
