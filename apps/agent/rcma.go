package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func getLANIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
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

	ip, _ := getLANIP()
	addr := fmt.Sprintf("%s:", ip)

	http.HandleFunc("/system-info", systemInfoHandler)

	ln, err := net.Listen("tcp", addr) // OS picks free port
	if err != nil {
		panic(err)
	}
	fmt.Println("Server running on", ln.Addr().String())
	http.Serve(ln, nil)
}
