@echo off
title VaultGuard - Advanced Secret Scanner
color 0B
echo.
echo   ██╗   ██╗ █████╗ ██╗   ██╗██╗  ████████╗
echo   ██║   ██║██╔══██╗██║   ██║██║  ╚══██╔══╝
echo   ██║   ██║███████║██║   ██║██║     ██║
echo   ╚██╗ ██╔╝██╔══██║██║   ██║██║     ██║
echo    ╚████╔╝ ██║  ██║╚██████╔╝███████╗██║
echo     ╚═══╝  ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝
echo    ██████╗ ██╗   ██╗ █████╗ ██████╗ ██████╗
echo   ██╔════╝ ██║   ██║██╔══██╗██╔══██╗██╔══██╗
echo   ██║  ███╗██║   ██║███████║██████╔╝██║  ██║
echo   ██║   ██║██║   ██║██╔══██║██╔══██╗██║  ██║
echo   ╚██████╔╝╚██████╔╝██║  ██║██║  ██║██████╔╝
echo    ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝
echo.
echo   Advanced Secret Scanner - 100%% Free ^& Local
echo   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo.

:: Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo   [ERROR] Go is not installed. Please install Go from https://go.dev
    pause
    exit /b 1
)

:: Build if binary doesn't exist
if not exist vaultguard.exe (
    echo   [1/2] Building VaultGuard engine...
    go build -o vaultguard.exe cmd/vaultguard/main.go
    if %ERRORLEVEL% neq 0 (
        echo   [ERROR] Build failed!
        pause
        exit /b 1
    )
    echo   [1/2] Build complete!
) else (
    echo   [1/2] Binary found, skipping build.
)

echo   [2/2] Starting dashboard on http://localhost:8080 ...
echo.
echo   ┌──────────────────────────────────────────┐
echo   │  Dashboard:  http://localhost:8080        │
echo   │  Press Ctrl+C to stop the server         │
echo   └──────────────────────────────────────────┘
echo.

:: Open browser automatically
start "" http://localhost:8080

:: Start server
vaultguard.exe serve -c pkg\scanner\rules.yaml

pause
