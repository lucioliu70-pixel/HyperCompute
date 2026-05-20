@echo off
setlocal

set SCRIPT_DIR=%~dp0
powershell -NoProfile -ExecutionPolicy Bypass -File "%SCRIPT_DIR%scripts\one_click_install_and_run_windows.ps1" %*
if errorlevel 1 (
  echo [HyperCompute] 一键安装并运行失败。
  exit /b 1
)

endlocal
