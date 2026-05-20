#!/usr/bin/env bash
set -e
psql ${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/hypercompute?sslmode=disable} <<SQL
create extension if not exists pgcrypto;
create table if not exists users(id serial primary key,email text unique,created_at timestamptz);
create table if not exists api_keys(id serial primary key,user_id int,key_hash text,created_at timestamptz);
create table if not exists nodes(id serial primary key,node_id text unique,owner_user_id int,base_url text,pool text,reputation_score float,health_score float,gpu_usage float,active_requests int,last_heartbeat_at timestamptz,status text,created_at timestamptz);
create table if not exists node_metrics_latest(node_id text primary key,gpu_usage float,vram_used_mb int,gpu_model text,updated_at timestamptz);
create table if not exists wallet_accounts(user_id int primary key,balance numeric,hold_balance numeric,updated_at timestamptz);
create table if not exists wallet_transactions(id serial primary key,user_id int,amount numeric,tx_type text,ref_id text,created_at timestamptz);
create table if not exists wallet_holds(id serial primary key,user_id int,ref_id text unique,amount numeric,status text,created_at timestamptz);
create table if not exists billing_records(id serial primary key,request_id text unique,user_id int,node_id text,model_name text,input_tokens int,output_tokens int,amount numeric,created_at timestamptz);
create table if not exists provider_earnings(id serial primary key,node_id text,billing_request_id text,amount numeric,status text,available_at timestamptz,created_at timestamptz);
create table if not exists withdrawal_requests(id serial primary key,node_id text,amount numeric,status text,created_at timestamptz);
create table if not exists risk_events(id serial primary key,event_type text,payload jsonb,created_at timestamptz);
create table if not exists verification_tasks(id serial primary key,node_id text,status text,created_at timestamptz);
create table if not exists verification_results(id serial primary key,task_id int,result text,created_at timestamptz);
create table if not exists models(id serial primary key,model_name text unique,status text,created_at timestamptz);
create table if not exists model_pricing(id serial primary key,model_name text unique,unit_price numeric,created_at timestamptz);
SQL
