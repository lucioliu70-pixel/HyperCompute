Write-Host "[HyperCompute] checking Docker Desktop..."
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) { throw "Docker 未安装" }
try { docker version | Out-Null } catch { throw "Docker Desktop 未启动" }
if (-not (Get-Command wsl -ErrorAction SilentlyContinue)) { Write-Warning "未检测到 WSL" }
make bootstrap
Start-Process "http://localhost:5173"
