package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	base := env("SCHEDULER_BASE_URL", "http://localhost:8081")
	nodeID := env("NODE_ID", "node-local")
	owner := env("OWNER_USER_ID", "1")
	nodeBase := env("NODE_BASE_URL", "http://localhost:8090")
	interval, _ := strconv.Atoi(env("HEARTBEAT_INTERVAL_SEC", "5"))
	post(base+"/nodes/register", map[string]any{"node_id": nodeID, "owner_user_id": atoi(owner), "base_url": nodeBase, "pool": env("NODE_POOL", "community"), "client_version": env("CLIENT_VERSION", "dev")})
	for {
		gpu, vr, model := probeGPU()
		post(base+"/nodes/heartbeat", map[string]any{"node_id": nodeID, "gpu_usage": gpu, "vram_used_mb": vr, "gpu_model": model, "runtime_online": checkRuntime(env("RUNTIME_BASE_URL", "http://localhost:8090"))})
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
func probeGPU() (float64, int, string) {
	out, err := exec.Command("bash", "-lc", "nvidia-smi --query-gpu=utilization.gpu,memory.used,name --format=csv,noheader,nounits | head -n1").Output()
	if err != nil {
		return 0, 0, "unknown"
	}
	p := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(p) < 3 {
		return 0, 0, "unknown"
	}
	g, _ := strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
	m, _ := strconv.Atoi(strings.TrimSpace(p[1]))
	return g, m, strings.TrimSpace(p[2])
}
func checkRuntime(base string) bool {
	resp, err := http.Get(base + "/v1/models")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}
func post(url string, p any) {
	b, _ := json.Marshal(p)
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(b))
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func atoi(s string) int { v, _ := strconv.Atoi(s); return v }
