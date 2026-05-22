@echo off
setlocal
title HyperCompute Launcher (Windows)

echo [HyperCompute] 正在检查环境...
where docker >nul 2>nul
if errorlevel 1 (
  echo [HyperCompute] 未检测到 Docker，尝试自动安装 Docker Desktop...
  where winget >nul 2>nul
  if errorlevel 1 (
    echo [HyperCompute] 未检测到 winget，无法自动安装 Docker。
    echo 请先安装 winget 或手动安装 Docker Desktop: https://www.docker.com/products/docker-desktop/
    pause
    exit /b 1
  )

  winget install -e --id Docker.DockerDesktop --accept-package-agreements --accept-source-agreements
  if errorlevel 1 (
    echo [HyperCompute] Docker 自动安装失败，请手动安装后重试。
    pause
    exit /b 1
  )

  echo [HyperCompute] Docker Desktop 安装完成，请先启动 Docker Desktop，然后重新运行本脚本。
  pause
  exit /b 0
)

echo [HyperCompute] Docker 命令已存在，检查 Docker Desktop 是否已启动...
docker version >nul 2>nul
if errorlevel 1 (
  echo [HyperCompute] Docker Desktop 尚未启动。
  echo 请先启动 Docker Desktop，待状态为 Running 后再重试。
  pause
  exit /b 1
)

where wsl >nul 2>nul || (echo [HyperCompute] 未检测到 WSL，建议安装 WSL2 以获得最佳体验)
where make >nul 2>nul
if errorlevel 1 (
  echo [HyperCompute] 未检测到 make 命令。
  echo 请在 Git Bash 中执行此脚本，或先安装 make（例如通过 MSYS2 / Chocolatey）。
  pause
  exit /b 1
)

make bootstrap
if errorlevel 1 (
  echo [HyperCompute] 启动失败，请根据上方错误信息排查。
  pause
  exit /b 1
)

start "" http://localhost:5173

echo [HyperCompute] 启动命令执行完成。
pause
endlocal
