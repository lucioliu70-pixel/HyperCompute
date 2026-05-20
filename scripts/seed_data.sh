#!/usr/bin/env bash
set -e
psql ${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/hypercompute?sslmode=disable} <<SQL
insert into users(id,email,created_at) values(1,'admin@hypercompute.local',now()) on conflict (id) do nothing;
insert into users(id,email,created_at) values(2,'node@hypercompute.local',now()) on conflict (id) do nothing;
insert into model_pricing(model_name,unit_price,created_at) values('Qwen/Qwen2.5-7B-Instruct',0.08,now()) on conflict (model_name) do nothing;
insert into models(model_name,status,created_at) values('Qwen/Qwen2.5-7B-Instruct','online',now()) on conflict (model_name) do nothing;
insert into wallet_accounts(user_id,balance,hold_balance,updated_at) values(1,100,0,now()) on conflict (user_id) do nothing;
SQL
