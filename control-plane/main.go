package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
	"strings"
)

func main() {
	db, _ := pgxpool.New(context.Background(), "postgres://postgres:postgres@postgres:5432/hypercompute?sslmode=disable")
	token := getenv("ADMIN_TOKEN", "admin-token")
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/admin/") {
				if r.Header.Get("Authorization") != "Bearer "+token {
					w.WriteHeader(401)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("cp_up 1")) })
	r.Get("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		var c, on int
		db.QueryRow(r.Context(), "select count(1) from users where contributor_status='APPROVED'").Scan(&c)
		db.QueryRow(r.Context(), "select count(1) from nodes where status='online' and now()-last_heartbeat_at<interval '30 second'").Scan(&on)
		json.NewEncoder(w).Encode(map[string]any{"contributors": c, "online_nodes": on})
	})
	r.Get("/admin/nodes", func(w http.ResponseWriter, r *http.Request) {
		rows, _ := db.Query(r.Context(), "select node_id,owner_user_id,base_url,status,runtime_online,gpu_usage,last_heartbeat_at from nodes order by id desc")
		defer rows.Close()
		a := []map[string]any{}
		for rows.Next() {
			var id, url, st string
			var owner int
			var rt bool
			var gpu float64
			var hb any
			rows.Scan(&id, &owner, &url, &st, &rt, &gpu, &hb)
			a = append(a, map[string]any{"node_id": id, "owner_user_id": owner, "base_url": url, "status": st, "runtime_online": rt, "gpu_usage": gpu, "last_heartbeat_at": hb})
		}
		json.NewEncoder(w).Encode(a)
	})
	r.Get("/admin/contributors", func(w http.ResponseWriter, r *http.Request) {
		rows, _ := db.Query(r.Context(), "select u.id,u.role,u.contributor_status,count(n.node_id) node_count,count(case when n.status='online' then 1 end) online_nodes,coalesce(sum(e.amount),0) earnings,coalesce(a.total_points,0) points from users u left join nodes n on n.owner_user_id=u.id left join provider_earnings e on e.node_id=n.node_id left join contributor_points_accounts a on a.user_id=u.id group by u.id,u.role,u.contributor_status,a.total_points order by u.id")
		defer rows.Close()
		o := []map[string]any{}
		for rows.Next() {
			var id, n, online int
			var role, st string
			var earn, pts float64
			rows.Scan(&id, &role, &st, &n, &online, &earn, &pts)
			o = append(o, map[string]any{"user_id": id, "role": role, "contributor_status": st, "node_count": n, "online_nodes": online, "earnings": earn, "points": pts})
		}
		json.NewEncoder(w).Encode(o)
	})
	r.Post("/admin/contributors/{user_id}/approve", func(w http.ResponseWriter, r *http.Request) {
		_, _ = db.Exec(r.Context(), "update users set contributor_status='APPROVED',role=case when role='user' then 'contributor' else 'both' end where id=$1", chi.URLParam(r, "user_id"))
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	r.Get("/contributor/me", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("user_id")
		var role, st string
		db.QueryRow(r.Context(), "select role,contributor_status from users where id=$1", uid).Scan(&role, &st)
		json.NewEncoder(w).Encode(map[string]any{"user_id": uid, "role": role, "contributor_status": st})
	})
	r.Get("/contributor/nodes", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("user_id")
		rows, _ := db.Query(r.Context(), "select node_id,status,runtime_online,last_heartbeat_at from nodes where owner_user_id=$1", uid)
		defer rows.Close()
		o := []map[string]any{}
		for rows.Next() {
			var id, st string
			var rt bool
			var hb any
			rows.Scan(&id, &st, &rt, &hb)
			o = append(o, map[string]any{"node_id": id, "status": st, "runtime_online": rt, "last_heartbeat_at": hb})
		}
		json.NewEncoder(w).Encode(o)
	})
	r.Get("/contributor/points", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("user_id")
		var t, a, f float64
		db.QueryRow(r.Context(), "select coalesce(total_points,0),coalesce(available_points,0),coalesce(frozen_points,0) from contributor_points_accounts where user_id=$1", uid).Scan(&t, &a, &f)
		json.NewEncoder(w).Encode(map[string]any{"total_points": t, "available_points": a, "frozen_points": f})
	})
	r.Get("/contributor/earnings", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("user_id")
		var amt float64
		db.QueryRow(r.Context(), "select coalesce(sum(e.amount),0) from provider_earnings e join nodes n on n.node_id=e.node_id where n.owner_user_id=$1", uid).Scan(&amt)
		json.NewEncoder(w).Encode(map[string]any{"earnings": amt})
	})
	http.ListenAndServe(":8085", r)
}
func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
