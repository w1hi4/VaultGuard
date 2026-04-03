# 🛡️ VaultGuard Nexus v3.1.0
### Advanced, Multi-Threaded Secret & PII Scanner

**VaultGuard Nexus** is an enterprise-grade, high-performance security auditing suite designed for the modern developer. It provides a "Zero-Configuration" experience for detecting hardcoded secrets, API keys, and sensitive PII across local and remote Git repositories.

---

## ⚡Core Features

- **🚀 Ultra-Scale Parallelism**: Leveraging a high-performance Go-based Worker Pool (8-16 workers) for blistering audit speeds on massive repositories.
- **🔍 Deep-History Intelligence**: Automatically traverses 100% of the commit history, unearthing secrets from deleted branches or historical commits.
- **💎 Nexus Defense UI**: A premium, WebGL-inspired dashboard providing real-time visual feedback and a "Digital Cascade" result reveal.
- **🌐 One-Click Remote Audits**: Simply paste any GitHub/GitLab URL. VaultGuard handles the cloning, analysis, and cleanup automatically.
- **🛡️ Intelligent Noise Reduction**: Automated path exclusion logic skips `node_modules`, `vendor`, `build`, and lock files for maximum performance.
- **📊 Professional Security Reporting**: Generates detailed, auditor-ready Markdown reports with one click from the dashboard.

---

## 🚀 Getting Started

### 📋 Prerequisites
- **Git**: [Download Git](https://git-scm.com/downloads) (Required for remote cloning)
- **Go 1.21+**: [Download Go](https://go.dev/dl/) (Required for building from source)

### 🛠️ Installation
```powershell
# Clone the repository
git clone https://github.com/w1hi4/VaultGuard.git
cd VaultGuard

# Build the high-performance binary
go build -o vaultguard.exe cmd/vaultguard/main.go
```

### 🔓 Quick Launch
The quickest way to start the **Nexus Dashboard** is via the provided batch script:
- **Windows**: Double-click `start.bat`
- **Linux/macOS**: `go run cmd/vaultguard/main.go serve -c pkg/scanner/rules.yaml`

---

## 🎨 Web Dashboard Guide
1. **Launch**: Start the server and navigate to `http://localhost:8080`.
2. **Audit**: Paste a **GitHub URL** or a **Local Directory Path** into the central "Nexus Link" field.
3. **Execute**: Click **AUDIT**. The system will immediately begin a deep, parallelized scan.
4. **Reveal**: Findings are presented in the "Digital Cascade" view, categorized by severity (CRITICAL, HIGH, MEDIUM, LOW).
5. **Report**: Click **Generate Report** to save a detailed security audit in Markdown format.

---

## 💻 CLI Master Guide
For power users and CI/CD integration, use the `vaultguard.exe` CLI.

| Command | Flag | Description | Example |
| :--- | :--- | :--- | :--- |
| **`scan`** | `-p`, `--path` | Target repository/path | `scan -p "https://github.com/user/repo"` |
| | `-e`, `--exclude` | Custom path patterns to skip | `scan -p "." -e "*.log,tmp/"` |
| | `--deep` | (Default: true) Deep scan history | `scan -p "." --deep=true` |
| | `--json` | Output findings in raw JSON | `scan -p "." --json > results.json` |
| **`serve`** | `-p`, `--port` | Web dashboard port | `serve --port 9090` |
| | `-c`, `--config` | Path to `rules.yaml` | `serve -c "pkg/scanner/rules.yaml"` |

---

## ⚙️ Customizing Detection Rules
Modify `pkg/scanner/rules.yaml` to add custom patterns or expand exclusions.

```yaml
rules:
  - id: my-custom-api-key
    description: "Detects internal secret headers"
    regex: "(?i)X-SECRET-KEY\\s*[:=]\\s*['\"]([0-9a-zA-Z]{32})['\"]"
    severity: "CRITICAL"

exclude_paths:
  - "node_modules/"
  - "dist/"
  - "*.log"
```

---

## 🛡️ License & Community
Distributed under the **MIT License**. Join the defense and help make the web more secure.

*Created by [w4hid] - Protecting your code, one commit at a time.*
