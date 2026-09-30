function Install-Winget {
    Write-Host "Winget not found." -ForegroundColor Yellow
    Write-Host "⌛ Installing winget..." -ForegroundColor Cyan
    try {
        $releasesUrl = "https://api.github.com/repos/microsoft/winget-cli/releases/latest"
        $latestRelease = Invoke-RestMethod -Uri $releasesUrl -ErrorAction Stop
        $asset = $latestRelease.assets | Where-Object { $_.name -like "*.msixbundle" } | Select-Object -First 1

        if ($asset) {
            $dest = Join-Path $env:TEMP $asset.name
            Write-Host "📥 Downloading $($asset.name)..." -ForegroundColor Cyan
            Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $dest -ErrorAction Stop

            Add-AppxPackage -Path $dest -ErrorAction Stop
            Remove-Item $dest -Force -ErrorAction SilentlyContinue
            
            $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
            
            if (Get-Command winget -ErrorAction SilentlyContinue) {
                Write-Host "✅ Winget has been successfully installed!" -ForegroundColor Green
                return $true
            }
        }
    } catch {
        Write-Host "⚠️️ Failed to install winget: $_" -ForegroundColor Yellow
    }
    return $false
}

function Install-Choco {
    Write-Host "⌛ Trying to install Chocolatey..." -ForegroundColor Cyan
    try {
        Set-ExecutionPolicy Bypass -Scope Process -Force
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
        
        $env:PATH = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

        if (Get-Command choco -ErrorAction SilentlyContinue) {
            Write-Host "✅ Chocolatey has been successfully installed!" -ForegroundColor Green
            return $true
        }
    } catch {
        Write-Host "⚠️ Error installing Chocolatey: $_" -ForegroundColor Red
    }
    return $false
}

function Install-ImageMagick {
    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
        Install-Winget | Out-Null
    }

    if (Get-Command winget -ErrorAction SilentlyContinue) {
        Write-Host "⌛ Installing ImageMagick via winget..." -ForegroundColor Cyan
        winget install ImageMagick.ImageMagick --accept-source-agreements --accept-package-agreements
        
        if ($LASTEXITCODE -eq 0) {
            return
        }
        Write-Host "⚠️ Installation error via winget. Fallback to Chocolatey..." -ForegroundColor Yellow
    }

    if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
        Install-Choco | Out-Null
    }

    if (Get-Command choco -ErrorAction SilentlyContinue) {
        Write-Host "⌛ Installing ImageMagick via Chocolatey..." -ForegroundColor Cyan
        choco install imagemagick -y
        if ($LASTEXITCODE -eq 0) {
            return
        }
    }

    Write-Host "❌ Error: Failed to install ImageMagick via both winget and choco." -ForegroundColor Red
    Write-Host "Install ImageMagick manually from https://imagemagick.org" -ForegroundColor Red
    exit 1
}

if (-not (Get-Command magick -ErrorAction SilentlyContinue) -and -not (Get-Command convert -ErrorAction SilentlyContinue)) {
    Write-Host "ImageMagick not found." -ForegroundColor Yellow
    Install-ImageMagick
} else {
    Write-Host "ImageMagick is already installed." -ForegroundColor Green
}

if (-not (Get-Command rsrc -ErrorAction SilentlyContinue)) {
    Write-Host "rsrc not found." -ForegroundColor Yellow
    Write-Host "⌛ Installing rsrc..." -ForegroundColor Cyan
    
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Host "Error: Go is not installed or not added to the PATH." -ForegroundColor Red
        exit 1
    }

    go install github.com/akavel/rsrc@latest

    $gopath = (go env GOPATH)
    $env:PATH += ";$gopath\bin"
    Write-Host "✅ rsrc has been successfully installed!" -ForegroundColor Green
} else {
    Write-Host "rsrc is already installed." -ForegroundColor Green
}