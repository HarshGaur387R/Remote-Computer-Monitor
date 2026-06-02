package agent

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo holds collected system metrics.
type SystemInfo struct {
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	TotalVRAMMB        uint64  `json:"total_vram_mb"`
	AvailableVRAMMB    uint64  `json:"available_vram_mb"`
	VRAMUsedPercent    float64 `json:"vram_used_percent"`
	TotalStorageGB     uint64  `json:"total_storage_gb"`
	AvailableStorageGB uint64  `json:"available_storage_gb"`
	StorageUsedPercent float64 `json:"storage_used_percent"`
}

type PingInfo struct {
	IsActive bool `json:"is_active"`
}

// CollectSystemInfo gathers CPU, memory, and disk metrics.
func CollectSystemInfo() (*SystemInfo, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("cpu.Percent: %w", err)
	}

	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("mem.VirtualMemory: %w", err)
	}

	d, err := disk.Usage("C:/")
	if err != nil {
		return nil, fmt.Errorf("disk.Usage: %w", err)
	}

	return &SystemInfo{
		CPUUsagePercent:    math.Round(percent[0]*100) / 100,
		TotalVRAMMB:        v.Total / 1024 / 1024,
		AvailableVRAMMB:    v.Available / 1024 / 1024,
		VRAMUsedPercent:    v.UsedPercent,
		TotalStorageGB:     d.Total / 1024 / 1024 / 1024,
		AvailableStorageGB: d.Free / 1024 / 1024 / 1024,
		StorageUsedPercent: math.Round(d.UsedPercent*100) / 100,
	}, nil
}

// AuthMiddleware rejects requests whose "auth" query param doesn't match token.
func AuthMiddleware(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("auth") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// SystemInfoHandler is the /system-info HTTP handler.
func SystemInfoHandler(w http.ResponseWriter, r *http.Request) {
	info, err := CollectSystemInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to collect metrics: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	responseData := &PingInfo{IsActive: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// BuildRouter creates and returns the HTTP mux for the agent.
func BuildRouter(authToken string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", AuthMiddleware(authToken, PingHandler))
	mux.HandleFunc("/system-info", AuthMiddleware(authToken, SystemInfoHandler))
	return mux
}

