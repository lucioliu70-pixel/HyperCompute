#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"
bash scripts/one_click_install_and_run.sh
read -r -p "按回车键退出..." _
