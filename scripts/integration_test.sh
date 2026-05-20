#!/usr/bin/env bash
set -e
bash scripts/health_check.sh
curl -s -XPOST localhost:8086/users -H 'content-type: application/json' -d '{"email":"user1@test.local"}' >/dev/null
curl -s -XPOST localhost:8086/users -H 'content-type: application/json' -d '{"email":"contrib@test.local"}' >/dev/null
curl -s -XPOST localhost:8086/contributors/apply -H 'content-type: application/json' -d '{"user_id":2,"display_name":"c1","contact_email":"c@test.local"}'|jq .ok|grep -q true
curl -s -XPOST localhost:8085/admin/contributors/2/approve -H 'Authorization: Bearer admin-token' -H 'content-type: application/json' -d '{}'|jq .ok|grep -q true
KEY=$(curl -s -XPOST localhost:8086/api-keys -H 'content-type: application/json' -d '{"user_id":1}'|jq -r .api_key)
curl -s -XPOST localhost:8082/wallets -H 'content-type: application/json' -d '{"user_id":1}' >/dev/null
curl -s -XPOST localhost:8082/topup -H 'content-type: application/json' -d '{"user_id":1,"amount":100}' >/dev/null
curl -s -XPOST localhost:8081/nodes/register -H 'content-type: application/json' -d '{"node_id":"node-001","owner_user_id":2,"base_url":"http://runtime-worker:8090","label":"n1","region":"us","tags":["4090"],"client_version":"1.0.0"}' >/dev/null
curl -s -XPOST localhost:8081/nodes/heartbeat -H 'content-type: application/json' -d '{"node_id":"node-001","gpu_usage":10,"vram_used_mb":1000,"gpu_model":"4090","runtime_online":true}' >/dev/null
curl -s -H 'Authorization: Bearer admin-token' localhost:8085/admin/nodes|jq length|grep -q '[1-9]'
curl -s -XPOST localhost:8080/v1/chat/completions -H "Authorization: Bearer $KEY" -H 'content-type: application/json' -d '{"model":"Qwen/Qwen2.5-7B-Instruct","stream":false,"messages":[{"role":"user","content":"hi"}]}'|jq . >/dev/null
curl -s -H 'Authorization: Bearer admin-token' localhost:8085/admin/provider-earnings|jq length|grep -q '[1-9]'
curl -s -XPOST localhost:8087/points/settle-available|jq .ok >/dev/null
echo HYPERCOMPUTE MVP PASS
