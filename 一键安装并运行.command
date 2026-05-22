#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"

if ! bash scripts/one_click_install_and_run.sh "$@"; then
  echo
  echo "❌ 一键脚本执行失败。请根据上方报错修复后重试。"
fi

read -r -p "按回车键退出..." _
