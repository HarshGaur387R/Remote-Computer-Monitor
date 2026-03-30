package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Config struct {
	LANIP     string `json:"lan_ip"`
	Port      int    `json:"port"`
	AuthToken string `json:"auth_token"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func authMiddleware(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("auth")
		if q != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func systemInfoHandler(w http.ResponseWriter, _ *http.Request) {

	percent, _ := cpu.Percent(time.Second, false)
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("C:/")

	data := map[string]any{
		"cpu_usage_percent": math.Round(percent[0]*100) / 100,

		"total_vram_mb":     v.Total / 1024 / 1024,
		"available_vram_mb": v.Available / 1024 / 1024,
		"vram_used_percent": v.UsedPercent,

		"total_storage_gb":     d.Total / 1024 / 1024 / 1024,
		"available_storage_gb": d.Free / 1024 / 1024 / 1024,
		"storage_used_percent": math.Round(d.UsedPercent*100) / 100,
	}
	json.NewEncoder(w).Encode(data)
}

func main() {
	fmt.Println("Hello from RCMA")

	// 1. Load config.json
	cfg, err := loadConfig(`C:\ProgramData\RCMAgent\config.json`)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Setup routes with authentication
	http.HandleFunc("/system-info", authMiddleware(cfg.AuthToken, systemInfoHandler))

	// 3. Bind to LAN IP + Port from config
	addr := fmt.Sprintf("%s:%d", cfg.LANIP, cfg.Port)
	fmt.Println("Agent running on http://" + addr)

	// 4. Start server
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
