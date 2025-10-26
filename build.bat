@echo off
REM Build script for Cultpedia on Windows

echo Building Cultpedia...
go build -o cultpedia.exe .\cmd

if %ERRORLEVEL% NEQ 0 (
    echo ✗ Build failed!
    exit /b 1
)

echo ✔ Build successful!
echo.
echo You can now run: .\cultpedia.exe
echo.
echo Usage:
echo   Interactive mode:  .\cultpedia.exe
echo   Commands:          .\cultpedia.exe help