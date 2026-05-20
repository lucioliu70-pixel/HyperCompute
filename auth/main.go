package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"hypercompute/shared"
)

type server struct{ db *pgxpool.Pool }

func main() {
	dsn := getenv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/hypercompute?sslmode=disable")
	db, _ := pgxpool.New(context.Background(), dsn)
	s := &server{db: db}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("auth_up 1")) })
	r.Post("/users", s.createUser)
	r.Get("/users/{id}", s.getUser)
	r.Post("/contributors/apply", s.applyContributor)
	r.Post("/contributors/approve", s.approveContributor)
	r.Post("/api-keys", s.createAPIKey)
	r.Post("/validate", s.validate)
	log.Fatal(http.ListenAndServe(":8086", r))
}
func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	var b struct{ Email, Role string }
	_ = json.NewDecoder(r.Body).Decode(&b)
	if b.Role == "" {
		b.Role = "user"
	}
	var id int
	_ = s.db.QueryRow(r.Context(), "insert into users(email,role,contributor_status,created_at) values($1,$2,'NONE',now()) on conflict(email) do update set email=excluded.email returning id", b.Email, b.Role).Scan(&id)
	json.NewEncoder(w).Encode(map[string]any{"id": id, "email": b.Email, "role": b.Role})
}
func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
	var id int
	var email, role, st string
	err := s.db.QueryRow(r.Context(), "select id,email,coalesce(role,'user'),coalesce(contributor_status,'NONE') from users where id=$1", chi.URLParam(r, "id")).Scan(&id, &email, &role, &st)
	if err != nil {
		shared.WriteError(w, 404, "user not found", "USER_NOT_FOUND")
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"id": id, "email": email, "role": role, "contributor_status": st})
}
func (s *server) applyContributor(w http.ResponseWriter, r *http.Request) {
	var b struct {
		UserID                    int    `json:"user_id"`
		DisplayName, ContactEmail string `json:"display_name" json:"contact_email"`
	}
	_ = json.NewDecoder(r.Body).Decode(&b)
	_, _ = s.db.Exec(r.Context(), "update users set contributor_status='PENDING' where id=$1", b.UserID)
	_, _ = s.db.Exec(r.Context(), "insert into contributor_profiles(user_id,display_name,contact_email,created_at,updated_at) values($1,$2,$3,now(),now()) on conflict(user_id) do update set display_name=excluded.display_name,contact_email=excluded.contact_email,updated_at=now()", b.UserID, b.DisplayName, b.ContactEmail)
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}
func (s *server) approveContributor(w http.ResponseWriter, r *http.Request) {
	var b struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&b)
	_, _ = s.db.Exec(r.Context(), "update users set contributor_status='APPROVED', role=case when role='user' then 'contributor' when role='admin' then 'admin' else 'both' end where id=$1", b.UserID)
	_, _ = s.db.Exec(r.Context(), "insert into contributor_points_accounts(user_id,total_points,available_points,frozen_points,updated_at) values($1,0,0,0,now()) on conflict(user_id) do nothing", b.UserID)
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}
func (s *server) createAPIKey(w http.ResponseWriter, r *http.Request) {
	var b struct {
		UserID int `json:"user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&b)
	raw := "hc_live_" + uuid.NewString()[:18]
	h := sha256.Sum256([]byte(raw))
	_, _ = s.db.Exec(r.Context(), "insert into api_keys(user_id,key_hash,created_at) values($1,$2,now())", b.UserID, hex.EncodeToString(h[:]))
	json.NewEncoder(w).Encode(map[string]any{"api_key": raw})
}
func (s *server) validate(w http.ResponseWriter, r *http.Request) {
	var b struct {
		APIKey string `json:"api_key"`
	}
	_ = json.NewDecoder(r.Body).Decode(&b)
	h := sha256.Sum256([]byte(b.APIKey))
	var uid int
	err := s.db.QueryRow(r.Context(), "select user_id from api_keys where key_hash=$1 order by id desc limit 1", hex.EncodeToString(h[:])).Scan(&uid)
	if err != nil {
		shared.WriteError(w, 401, "invalid api key", "INVALID_API_KEY")
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"valid": true, "user_id": uid})
}
func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
