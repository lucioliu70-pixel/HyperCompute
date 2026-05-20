package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func main() {
	db, _ := pgxpool.New(context.Background(), "postgres://postgres:postgres@postgres:5432/hypercompute?sslmode=disable")
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("settlement_up 1")) })
	r.Post("/earnings", func(w http.ResponseWriter, r *http.Request) {
		var b map[string]any
		json.NewDecoder(r.Body).Decode(&b)
		amt := b["amount"].(float64) * 0.65
		_, _ = db.Exec(r.Context(), "insert into provider_earnings(node_id,billing_request_id,amount,status,available_at,created_at) values($1,$2,$3,'PENDING',now()+interval '7 day',now())", b["node_id"], b["request_id"], amt)
		_, _ = db.Exec(r.Context(), "insert into contributor_points_ledger(user_id,node_id,source_type,source_ref_id,delta_points,status,created_at) select owner_user_id,$1,'EARNING',$2,$3,'PENDING',now() from nodes where node_id=$1", b["node_id"], b["request_id"], amt*1000)
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	r.Post("/points/settle-available", func(w http.ResponseWriter, r *http.Request) {
		_, _ = db.Exec(r.Context(), "update contributor_points_ledger l set status='AVAILABLE' where status='PENDING' and exists (select 1 from provider_earnings e where e.billing_request_id=l.source_ref_id and e.available_at<=now())")
		_, _ = db.Exec(r.Context(), "insert into contributor_points_accounts(user_id,total_points,available_points,frozen_points,updated_at) select user_id,sum(delta_points),sum(case when status='AVAILABLE' then delta_points else 0 end),sum(case when status='PENDING' then delta_points else 0 end),now() from contributor_points_ledger group by user_id on conflict(user_id) do update set total_points=excluded.total_points,available_points=excluded.available_points,frozen_points=excluded.frozen_points,updated_at=now()")
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	http.ListenAndServe(":8087", r)
}
