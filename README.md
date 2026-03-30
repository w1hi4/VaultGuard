# 🛡️ VaultGuard: Advanced Secret & PII Scanner

**VaultGuard** is a high-performance, local-first security tool designed to detect hardcoded secrets, API keys, and PII (Personally Identifiable Information) across your entire Git history.

![VaultGuard Dashboard Screenshot](https://via.placeholder.com/1200x600/1e293b/ffffff?text=VaultGuard+Premium+Dashboard)

## ✨ Features

- **🚀 Blistering Speed**: Built in Go for rapid analysis of even the largest repositories.
- **🔍 Deep History Scanning**: Automatically analyzes the full commit chain, uncovering secrets buried years ago.
- **🌐 Remote Git Support**: Paste any Git URL (GitHub/GitLab/etc.) to automatically clone and scan.
- **💎 Premium UI**: A modern, glassmorphic dashboard with real-time feedback and server-side analysis.
- **📊 Professional Reporting**: One-click generation of detailed Markdown audit reports for compliance and security.
- **🛠️ Extensible Rules**: Add your own detection patterns via a simple `rules.yaml` file.

## 🚀 Quick Start (Local)

### Prerequisites
- [Go 1.21+](https://go.dev/dl/) installed.
- [Git](https://git-scm.com/downloads) installed.

### Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/vaultguard.git
cd vaultguard

# Build the engine
go build -o vaultguard.exe cmd/vaultguard/main.go
```

### Running the Dashboard
The easiest way to start is using the provided launcher:
- **Windows**: Double-click `start.bat`
- **Linux/macOS**: `go run cmd/vaultguard/main.go serve -c pkg/scanner/rules.yaml`

Then open `http://localhost:8080` in your browser.

## 💻 CLI Usage

Scan a local project or remote repository directly from the terminal:

```bash
# Scan a local folder
./vaultguard.exe scan -p "C:/projects/myapp"

# Deep scan a remote GitHub repository
./vaultguard.exe scan -p "https://github.com/user/repo.git" --deep

# Export results as JSON
./vaultguard.exe scan -p "https://github.com/user/repo.git" --json > audit.json
```

## 🛡️ License
Distributed under the **MIT License**. See `LICENSE` for more information.

---
*Created by [Your Name/Handle] - Protecting your code, one commit at a time.*
