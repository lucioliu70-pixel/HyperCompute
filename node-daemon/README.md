# Node Daemon

## 环境变量
SCHEDULER_BASE_URL, NODE_ID, OWNER_USER_ID, NODE_BASE_URL, NODE_POOL, CLIENT_VERSION, HEARTBEAT_INTERVAL_SEC, RUNTIME_BASE_URL

## Linux
`go run ./node-daemon`

## Windows (PowerShell)
`$env:SCHEDULER_BASE_URL='http://localhost:8081'; go run .\node-daemon`
