@echo off
setlocal

echo [HyperCompute] checking Docker Desktop...
docker version >nul 2>nul || (
  echo Docker Desktop 未启动，请先启动后重试
  pause
  exit /b 1
)
where wsl >nul 2>nul || (echo 未检测到 WSL，建议安装 WSL2 以获得最佳体验)

make bootstrap
if errorlevel 1 (
  echo [HyperCompute] 启动失败，请根据上方错误信息排查。
  pause
  exit /b 1
)

start http://localhost:5173

echo [HyperCompute] 启动命令执行完成。
pause
endlocal
