package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, _ := pgxpool.New(context.Background(), getenv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/hypercompute?sslmode=disable"))
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("scheduler_up 1")) })
	r.Post("/nodes/register", func(w http.ResponseWriter, r *http.Request) {
		var b map[string]any
		_ = json.NewDecoder(r.Body).Decode(&b)
		_, _ = db.Exec(r.Context(), "insert into nodes(node_id,owner_user_id,base_url,pool,label,region,tags,client_version,reputation_score,health_score,gpu_usage,active_requests,last_heartbeat_at,status,runtime_online,created_at) values($1,$2,$3,$4,$5,$6,$7::jsonb,$8,100,100,0,0,now(),'online',false,now()) on conflict(node_id) do update set base_url=excluded.base_url,pool=excluded.pool,label=excluded.label,region=excluded.region,tags=excluded.tags,client_version=excluded.client_version,status='online',last_heartbeat_at=now()", b["node_id"], int(b["owner_user_id"].(float64)), b["base_url"], coalesce(b["pool"], "community"), b["label"], b["region"], toJSON(b["tags"]), b["client_version"])
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	r.Post("/nodes/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		var b map[string]any
		_ = json.NewDecoder(r.Body).Decode(&b)
		_, _ = db.Exec(r.Context(), "update nodes set last_heartbeat_at=now(),status='online',gpu_usage=$2,runtime_online=$3 where node_id=$1", b["node_id"], b["gpu_usage"], b["runtime_online"])
		_, _ = db.Exec(r.Context(), "insert into node_metrics_latest(node_id,gpu_usage,vram_used_mb,gpu_model,updated_at) values($1,$2,$3,$4,now()) on conflict(node_id) do update set gpu_usage=excluded.gpu_usage,vram_used_mb=excluded.vram_used_mb,gpu_model=excluded.gpu_model,updated_at=now()", b["node_id"], b["gpu_usage"], b["vram_used_mb"], b["gpu_model"])
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	r.Get("/nodes", func(w http.ResponseWriter, r *http.Request) {
		q := "select node_id,owner_user_id,base_url,status,runtime_online,gpu_usage,last_heartbeat_at from nodes where 1=1"
		args := []any{}
		if v := r.URL.Query().Get("owner_user_id"); v != "" {
			q += " and owner_user_id=$1"
			id, _ := strconv.Atoi(v)
			args = append(args, id)
		}
		rows, _ := db.Query(r.Context(), q, args...)
		defer rows.Close()
		out := []map[string]any{}
		for rows.Next() {
			var id, url, st string
			var owner int
			var rt bool
			var gpu float64
			var hb time.Time
			rows.Scan(&id, &owner, &url, &st, &rt, &gpu, &hb)
			if time.Since(hb) > 30*time.Second {
				st = "offline"
			}
			out = append(out, map[string]any{"node_id": id, "owner_user_id": owner, "base_url": url, "status": st, "runtime_online": rt, "gpu_usage": gpu, "last_heartbeat_at": hb})
		}
		json.NewEncoder(w).Encode(out)
	})
	r.Post("/nodes/{node_id}/set-status", func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Status string `json:"status"`
		}
		_ = json.NewDecoder(r.Body).Decode(&b)
		_, _ = db.Exec(r.Context(), "update nodes set status=$2 where node_id=$1", chi.URLParam(r, "node_id"), b.Status)
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	r.Post("/schedule", func(w http.ResponseWriter, r *http.Request) {
		var b map[string]any
		_ = json.NewDecoder(r.Body).Decode(&b)
		uid := int(b["user_id"].(float64))
		pref, _ := b["preferred_node_id"].(string)
		q := "select node_id,base_url from nodes where owner_user_id<>$1 and status='online' and now()-last_heartbeat_at < interval '30 second'"
		args := []any{uid}
		if pref != "" {
			q += " order by case when node_id=$2 then 0 else 1 end,reputation_score desc,health_score desc,gpu_usage asc,active_requests asc limit 1"
			args = append(args, pref)
		} else {
			q += " order by reputation_score desc,health_score desc,gpu_usage asc,active_requests asc limit 1"
		}
		var id, url string
		_ = db.QueryRow(r.Context(), q, args...).Scan(&id, &url)
		json.NewEncoder(w).Encode(map[string]any{"node_id": id, "base_url": url, "ts": time.Now()})
	})
	http.ListenAndServe(":8081", r)
}
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func coalesce(v any, d string) any {
	if v == nil || v == "" {
		return d
	}
	return v
}
func toJSON(v any) string {
	b, _ := json.Marshal(v)
	if string(b) == "null" {
		return "[]"
	}
	return string(b)
}
