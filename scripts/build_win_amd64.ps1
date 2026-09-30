$ErrorActionPreference = "Stop"

$ProjectDir = Get-Location
$BuildDir   = Join-Path $ProjectDir "build\win\amd64"
$IconsDir   = Join-Path $ProjectDir "icons"

$SrcImg = Join-Path $IconsDir "icon_win.png"
if (-not (Test-Path $SrcImg)) {
    $SrcImg = Join-Path $IconsDir "icon_mac.png"
}

$IcoPath  = Join-Path $IconsDir "icon.ico"
$AppDir = Join-Path $ProjectDir "cmd\app"
$SysoPath = Join-Path $AppDir "\rsrc_windows_amd64.syso"

Write-Host "⌛ Preparing the icon..." -ForegroundColor Cyan

$env:PATH = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

$magickExe = (Get-Command magick.exe -ErrorAction SilentlyContinue).Source

if (-not $magickExe) {
    $possiblePaths = Get-ChildItem -Path "$env:ProgramFiles\ImageMagick*" -Filter "magick.exe" -Recurse -ErrorAction SilentlyContinue
    if ($possiblePaths) {
        $magickExe = $possiblePaths[0].FullName
    }
}

if ($magickExe -and (Test-Path $SrcImg)) {
    Write-Host "⌛ Converting $SrcImg to $IcoPath via ImageMagick..." -ForegroundColor Cyan
    & $magickExe "$SrcImg" -define icon:auto-resize="256,128,64,48,32,16" "$IcoPath"
} else {
    Write-Host "⚠️ ImageMagick (magick.exe) not found. Checking if $IcoPath already exists..." -ForegroundColor Yellow
}

if (-not $IcoPath -or -not (Test-Path $IcoPath)) {
    Write-Host "❌ Error: Failed to create or find icon at: $IcoPath" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Icon prepared successfully: $IcoPath" -ForegroundColor Green

if (Get-Command rsrc -ErrorAction SilentlyContinue) {
    Write-Host "⌛ Generating Windows resource file (.syso)..." -ForegroundColor Cyan
    rsrc -ico "$IcoPath" -arch arm64 -o "$SysoPath"
} else {
    Write-Host "⚠️ Warning: 'rsrc' utility not found. Icon will not be embedded into the .exe binary." -ForegroundColor Yellow
}

Write-Host "🔨 Building application..." -ForegroundColor Cyan
if (-not (Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
}

$env:CGO_ENABLED = "1"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

$ExePath = Join-Path $BuildDir "WordToExcel.exe"

try {
    cd "$AppDir"
    go build -ldflags="-H windowsgui" -x -o "$ExePath" -buildvcs=false .
    Write-Host "✅ Build finished successfully: $ExePath" -ForegroundColor Green
    cd "$ProjectDir"
} finally {
    if (Test-Path $SysoPath) {
        Remove-Item "$SysoPath" -Force -ErrorAction SilentlyContinue
    }
}

Write-Host "⌛ Creating release zip archive..." -ForegroundColor Cyan
$ZipPath = Join-Path $BuildDir "WordToExcel_win_amd64.zip"

if (Test-Path $ZipPath) {
    Remove-Item $ZipPath -Force
}

Compress-Archive -Path "$ExePath" -DestinationPath "$ZipPath" -Force
Write-Host "✅ Created ZIP package: $ZipPath" -ForegroundColor Green