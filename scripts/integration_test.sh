#!/usr/bin/env bash
set -e
bash scripts/health_check.sh
curl -s -XPOST localhost:8086/users -H 'content-type: application/json' -d '{"email":"user1@test.local"}' >/dev/null
KEY=$(curl -s -XPOST localhost:8086/api-keys -H 'content-type: application/json' -d '{"user_id":1}'|jq -r .api_key)
curl -s -XPOST localhost:8082/wallets -H 'content-type: application/json' -d '{"user_id":1}' >/dev/null
curl -s -XPOST localhost:8082/topup -H 'content-type: application/json' -d '{"user_id":1,"amount":100}' >/dev/null
curl -s localhost:8088/models|jq . >/dev/null
curl -s -XPOST localhost:8081/nodes/register -H 'content-type: application/json' -d '{"node_id":"node-001","owner_user_id":2,"base_url":"http://runtime-worker:8090"}' >/dev/null
curl -s -XPOST localhost:8081/nodes/heartbeat -H 'content-type: application/json' -d '{"node_id":"node-001"}' >/dev/null
curl -s -XPOST localhost:8080/v1/chat/completions -H "Authorization: Bearer $KEY" -H 'content-type: application/json' -d '{"model":"Qwen/Qwen2.5-7B-Instruct","stream":false,"messages":[{"role":"user","content":"hi"}]}'|jq . >/dev/null
B=$(curl -s localhost:8082/wallets/1|jq -r .balance); awk "BEGIN{exit !($B<200)}"
curl -s -H 'Authorization: Bearer admin-token' localhost:8085/admin/billing-records|jq length|grep -q '[1-9]'
curl -s -H 'Authorization: Bearer admin-token' localhost:8085/admin/provider-earnings|jq length|grep -q '[1-9]'
curl -s -H 'Authorization: Bearer admin-token' localhost:8085/admin/dashboard|jq . >/dev/null
curl -sN -XPOST localhost:8080/v1/chat/completions -H "Authorization: Bearer $KEY" -H 'content-type: application/json' -d '{"model":"Qwen/Qwen2.5-7B-Instruct","stream":true,"messages":[{"role":"user","content":"hi"}]}'|head -n 1|grep -q data:
echo HYPERCOMPUTE MVP PASS
