# HyperCompute MVP

## 一键启动
- Linux/macOS: `make bootstrap`
- Windows: `start_windows.bat` 或 `start_windows.ps1`

## 使用者流程
1. 创建用户 `/users`
2. 创建 API Key `/api-keys`
3. 调用 `/v1/chat/completions`

## 贡献者流程
1. 申请贡献者 `POST /contributors/apply`
2. 管理员审批 `POST /admin/contributors/{user_id}/approve`
3. 下载并运行 node-daemon（见 `node-daemon/README.md`）
4. 管理后台查看 `/admin/nodes`、`/admin/contributors`
5. 收益和积分：`/contributor/earnings`、`/contributor/points`

## 测试
`bash scripts/integration_test.sh`
