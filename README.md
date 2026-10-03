# WordToExcel

<p align="center">
  <img src="./icons/icon_mac.png" alt="WordToExcel logo" width="250" />
</p>

WordToExcel is a Go desktop application for converting tables from DOCX files into Excel (.xlsx). It is useful when a Word document contains one or more tables that need to be quickly transferred into spreadsheet format for further analysis, editing, or import.

## Features

- Select one or more DOCX files
- Extract tables from Word documents
- Automatically detect the main table in each file
- Merge similar tables and move additional tables to separate sheets
- Export results to Excel while preserving table structure and data
- Cross-platform GUI for macOS, Windows, and Linux

## Tech Stack

- Go
- Gio UI (`gioui.org`) for the interface
- `github.com/nguyenthenguyen/docx` for DOCX parsing
- `github.com/xuri/excelize/v2` for Excel generation

## Requirements

- Go 1.26.3 or newer
- On Linux, Windows, or macOS, GUI dependencies may need to be installed according to the platform
- Platform-specific build scripts are available in the `scripts/` directory

## Quick Start

1. Clone the repository:

   ```bash
   git clone <repo-url>
   cd WordToExcel
   ```

2. Download dependencies:

   ```bash
   go mod download
   ```

3. Run the application:

   ```bash
   go run ./cmd/app
   ```

## Build

### Local build

```bash
go build -o wordtoexcel ./cmd/app
```

### Platform builds

The project includes build scripts for common targets:

```bash
./scripts/build_mac_amd64.sh
./scripts/build_mac_arm64.sh
./scripts/build_linux_amd64.sh
./scripts/build_linux_arm64.sh
./scripts/build_win_amd64.sh
./scripts/build_win_arm64.sh
```

PowerShell build scripts are also available for Windows:

```powershell
.\scripts\build_win_amd64.ps1
.\scripts\build_win_arm64.ps1
```

## Usage

1. Start the application.
2. Click `Select Files` and choose the DOCX documents.
3. After selecting files, click `Convert`.
4. The application extracts tables and generates a temporary Excel file.
5. Once processing is complete, save the result to the desired location.

## [LICENSE](./LICENSE)

