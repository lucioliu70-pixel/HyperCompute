@echo off
echo [HyperCompute] checking Docker Desktop...
docker version >nul 2>nul || (echo Docker Desktop 未启动，请先启动后重试 & exit /b 1)
where wsl >nul 2>nul || (echo 未检测到 WSL，建议安装 WSL2 以获得最佳体验)
make bootstrap || exit /b 1
start http://localhost:5173
