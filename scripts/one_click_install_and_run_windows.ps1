Param(
  [switch]$Lite
)

$ErrorActionPreference = 'Stop'

function Log([string]$Message) {
  Write-Host "`n[HyperCompute] $Message"
}

function Ensure-Command([string]$Name) {
  return [bool](Get-Command $Name -ErrorAction SilentlyContinue)
}

function Ensure-Docker {
  if (-not (Ensure-Command "docker")) {
    throw "未检测到 Docker。请先安装并启动 Docker Desktop：https://docs.docker.com/desktop/setup/install/windows-install/"
  }

  try {
    docker info | Out-Null
  }
  catch {
    throw "Docker Desktop 未启动，请先启动后重试。"
  }
}

function Ensure-Compose {
  try {
    docker compose version | Out-Null
  }
  catch {
    throw "未检测到 docker compose，请升级 Docker Desktop。"
  }
}

function Invoke-Step([string]$Label, [scriptblock]$Action) {
  Log $Label
  & $Action
}

$RootDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$ComposeFile = Join-Path $RootDir "deploy/docker-compose.yml"
$LiteComposeFile = Join-Path $RootDir "deploy/docker-compose.lite.yml"
$Mode = if ($Lite) { "lite" } else { "full" }
$TargetComposeFile = if ($Lite) { $LiteComposeFile } else { $ComposeFile }

Set-Location $RootDir

Log "开始一键安装并运行 HyperCompute (Windows)"
Ensure-Docker
Ensure-Compose

Invoke-Step "执行项目初始化（依赖拉取、构建等）" { bash "$RootDir/deploy/bootstrap.sh" }
Invoke-Step "启动系统服务（模式：$Mode）" { docker compose -f "$TargetComposeFile" up -d --build }
Invoke-Step "等待数据库启动" { bash "$RootDir/scripts/wait_for_postgres.sh" }
Invoke-Step "执行数据库迁移" { bash "$RootDir/scripts/migrate_all.sh" }
Invoke-Step "执行初始化数据" { bash "$RootDir/scripts/seed_data.sh" }
Invoke-Step "执行健康检查" { bash "$RootDir/scripts/health_check.sh" }

Log "✅ 环境安装与系统部署完成。"
Log "查看日志: make logs"
Log "停止系统: make down"
