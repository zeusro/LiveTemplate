// Part 3: 简易服务器状态监控看板 - Go 后端
// 提供 CPU/内存模拟数据、SSE 推送、可动态修改推送频率

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	CPU  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
	Time int64   `json:"time"`
}

var (
	// 推送间隔（秒），通过 API 动态修改
	intervalSec atomic.Int64
	// 用于 SSE 广播：新连接订阅 channel
	subscribers   = make(map[chan []byte]struct{})
	subscribersMu sync.RWMutex
)

func main() {
	intervalSec.Store(1) // 默认 1 秒

	mux := http.NewServeMux()

	// 一次性的指标 API（随机 CPU/内存）
	mux.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(genMetrics())
	})

	// SSE 流：按当前间隔持续推送
	mux.HandleFunc("/api/stream", handleSSE)

	// 修改推送频率，例如 PUT /api/interval?seconds=5 或 POST body: {"seconds": 5}
	mux.HandleFunc("/api/interval", handleInterval)

	addr := ":8080"
	log.Printf("backend listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func genMetrics() Metrics {
	return Metrics{
		CPU:  round(30+rand.Float64()*50, 2),
		Mem:  round(40+rand.Float64()*40, 2),
		Time: time.Now().UnixMilli(),
	}
}

func round(v float64, n int) float64 {
	pow := 1.0
	for i := 0; i < n; i++ {
		pow *= 10
	}
	return float64(int(v*pow+0.5)) / pow
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := make(chan []byte, 8)
	subscribersMu.Lock()
	subscribers[ch] = struct{}{}
	subscribersMu.Unlock()
	defer func() {
		subscribersMu.Lock()
		delete(subscribers, ch)
		subscribersMu.Unlock()
		close(ch)
	}()

	// 通知客户端连接成功
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			if _, err := w.Write([]byte("data: " + string(data) + "\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// 启动一个全局 ticker，按当前间隔向所有 SSE 连接广播
func init() {
	go func() {
		ticker := time.NewTicker(time.Duration(intervalSec.Load()) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			interval := intervalSec.Load()
			ticker.Reset(time.Duration(interval) * time.Second)
			data, _ := json.Marshal(genMetrics())
			subscribersMu.RLock()
			for sub := range subscribers {
				select {
				case sub <- data:
				default:
				}
			}
			subscribersMu.RUnlock()
		}
	}()
}

func handleInterval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var sec int64 = 1
	if s := r.URL.Query().Get("seconds"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v >= 1 && v <= 60 {
			sec = v
		}
	} else if r.Method == http.MethodPost || r.Method == http.MethodPut {
		var body struct {
			Seconds int64 `json:"seconds"`
		}
		if json.NewDecoder(r.Body).Decode(&body) == nil && body.Seconds >= 1 && body.Seconds <= 60 {
			sec = body.Seconds
		}
	}

	intervalSec.Store(sec)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"seconds": sec})
}
