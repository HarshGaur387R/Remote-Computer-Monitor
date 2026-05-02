package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "net/http"
    "os"
    "time"

    "golang.org/x/sys/windows/svc"
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
        "cpu_usage_percent":    math.Round(percent[0]*100) / 100,
        "total_vram_mb":        v.Total / 1024 / 1024,
        "available_vram_mb":    v.Available / 1024 / 1024,
        "vram_used_percent":    v.UsedPercent,
        "total_storage_gb":     d.Total / 1024 / 1024 / 1024,
        "available_storage_gb": d.Free / 1024 / 1024 / 1024,
        "storage_used_percent": math.Round(d.UsedPercent*100) / 100,
    }
    json.NewEncoder(w).Encode(data)
}

// ---- Service wrapper ----
type rcmaService struct {
    server *http.Server
}

func (m *rcmaService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
    const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
    s <- svc.Status{State: svc.StartPending}
    s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

    // Start agent in goroutine
    go m.runAgent()

    for c := range r {
        switch c.Cmd {
        case svc.Stop, svc.Shutdown:
            s <- svc.Status{State: svc.StopPending}
            // Graceful shutdown
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            if err := m.server.Shutdown(ctx); err != nil {
                log.Printf("Error shutting down server: %v", err)
            }
            return false, 0
        }
    }
    return false, 0
}

func (m *rcmaService) runAgent() {
    cfg, err := loadConfig(`C:\ProgramData\RCMA\config.json`)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/system-info", authMiddleware(cfg.AuthToken, systemInfoHandler))

    addr := fmt.Sprintf("%s:%d", cfg.LANIP, cfg.Port)
    m.server = &http.Server{
        Addr:    addr,
        Handler: mux,
    }

    log.Println("Agent running on http://" + addr)
    if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Server failed: %v", err)
    }
}

func main() {
    isService, err := svc.IsWindowsService()
    if err != nil {
        log.Fatal(err)
    }
    if isService {
        svc.Run("rcma", &rcmaService{})
    } else {
        // Console mode for debugging
        s := &rcmaService{}
        s.runAgent()
    }
}
