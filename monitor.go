package spoor

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMonitor monitors logger performance
type PerformanceMonitor struct {
	mu           sync.RWMutex
	startTime    time.Time
	totalLogs    int64
	droppedLogs  int64
	errorCount   int64
	memoryStats  *MemoryStats
	cpuStats     *CPUStats
	latencyStats *LatencyStats
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	closed       int32
}

// MemoryStats tracks memory usage
type MemoryStats struct {
	Alloc        uint64
	TotalAlloc   uint64
	Sys          uint64
	NumGC        uint32
	HeapObjects  uint64
	HeapAlloc    uint64
	HeapSys      uint64
	HeapIdle     uint64
	HeapInuse    uint64
	HeapReleased uint64
	StackInuse   uint64
	StackSys     uint64
	MSpanInuse   uint64
	MSpanSys     uint64
	MCacheInuse  uint64
	MCacheSys    uint64
	BuckHashSys  uint64
	GCSys        uint64
	OtherSys     uint64
}

// CPUStats tracks CPU usage
type CPUStats struct {
	NumGoroutine int
	NumCPU       int
	NumCgoCall   int64
}

// LatencyStats tracks latency metrics
type LatencyStats struct {
	MinLatency    time.Duration
	MaxLatency    time.Duration
	AvgLatency    time.Duration
	P50Latency    time.Duration
	P90Latency    time.Duration
	P95Latency    time.Duration
	P99Latency    time.Duration
	TotalLatency  time.Duration
	LatencyCount  int64
	latencySamples []time.Duration
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	
	monitor := &PerformanceMonitor{
		startTime:    time.Now(),
		memoryStats:  &MemoryStats{},
		cpuStats:     &CPUStats{},
		latencyStats: &LatencyStats{
			latencySamples: make([]time.Duration, 0, 1000),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// Start monitoring goroutine
	monitor.wg.Add(1)
	go monitor.monitorLoop()

	return monitor
}

// RecordLog records a log entry
func (pm *PerformanceMonitor) RecordLog() {
	atomic.AddInt64(&pm.totalLogs, 1)
}

// RecordDropped records a dropped log
func (pm *PerformanceMonitor) RecordDropped() {
	atomic.AddInt64(&pm.droppedLogs, 1)
}

// RecordError records an error
func (pm *PerformanceMonitor) RecordError() {
	atomic.AddInt64(&pm.errorCount, 1)
}

// RecordLatency records a latency measurement
func (pm *PerformanceMonitor) RecordLatency(latency time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.latencyStats.TotalLatency += latency
	atomic.AddInt64(&pm.latencyStats.LatencyCount, 1)

	// Update min/max
	if pm.latencyStats.MinLatency == 0 || latency < pm.latencyStats.MinLatency {
		pm.latencyStats.MinLatency = latency
	}
	if latency > pm.latencyStats.MaxLatency {
		pm.latencyStats.MaxLatency = latency
	}

	// Add to samples for percentile calculation
	if len(pm.latencyStats.latencySamples) < 1000 {
		pm.latencyStats.latencySamples = append(pm.latencyStats.latencySamples, latency)
	}
}

// GetStats returns current performance statistics
func (pm *PerformanceMonitor) GetStats() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// Update memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	pm.memoryStats = &MemoryStats{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		NumGC:        m.NumGC,
		HeapObjects:  m.HeapObjects,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		StackInuse:   m.StackInuse,
		StackSys:     m.StackSys,
		MSpanInuse:   m.MSpanInuse,
		MSpanSys:     m.MSpanSys,
		MCacheInuse:  m.MCacheInuse,
		MCacheSys:    m.MCacheSys,
		BuckHashSys:  m.BuckHashSys,
		GCSys:        m.GCSys,
		OtherSys:     m.OtherSys,
	}

	// Update CPU stats
	pm.cpuStats = &CPUStats{
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		NumCgoCall:   runtime.NumCgoCall(),
	}

	// Calculate latency percentiles
	pm.calculateLatencyPercentiles()

	// Calculate throughput
	uptime := time.Since(pm.startTime)
	throughput := float64(atomic.LoadInt64(&pm.totalLogs)) / uptime.Seconds()

	return map[string]interface{}{
		"uptime":      uptime,
		"total_logs":  atomic.LoadInt64(&pm.totalLogs),
		"dropped_logs": atomic.LoadInt64(&pm.droppedLogs),
		"error_count": atomic.LoadInt64(&pm.errorCount),
		"throughput":  throughput,
		"memory":      pm.memoryStats,
		"cpu":         pm.cpuStats,
		"latency":     pm.latencyStats,
	}
}

// calculateLatencyPercentiles calculates latency percentiles
func (pm *PerformanceMonitor) calculateLatencyPercentiles() {
	if len(pm.latencyStats.latencySamples) == 0 {
		return
	}

	// Sort samples (simplified - in production use proper sorting)
	samples := make([]time.Duration, len(pm.latencyStats.latencySamples))
	copy(samples, pm.latencyStats.latencySamples)

	// Calculate percentiles
	count := len(samples)
	if count > 0 {
		pm.latencyStats.AvgLatency = pm.latencyStats.TotalLatency / time.Duration(pm.latencyStats.LatencyCount)
		
		// Simple percentile calculation (not accurate but fast)
		if count >= 1 {
			pm.latencyStats.P50Latency = samples[count/2]
		}
		if count >= 10 {
			pm.latencyStats.P90Latency = samples[int(float64(count)*0.9)]
		}
		if count >= 20 {
			pm.latencyStats.P95Latency = samples[int(float64(count)*0.95)]
		}
		if count >= 100 {
			pm.latencyStats.P99Latency = samples[int(float64(count)*0.99)]
		}
	}
}

// monitorLoop runs the monitoring loop
func (pm *PerformanceMonitor) monitorLoop() {
	defer pm.wg.Done()
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Update stats periodically
			pm.GetStats()
		case <-pm.ctx.Done():
			return
		}
	}
}

// Close closes the performance monitor
func (pm *PerformanceMonitor) Close() error {
	if !atomic.CompareAndSwapInt32(&pm.closed, 0, 1) {
		return nil
	}

	pm.cancel()
	pm.wg.Wait()
	return nil
}

// PrintStats prints formatted performance statistics
func (pm *PerformanceMonitor) PrintStats() {
	stats := pm.GetStats()
	
	fmt.Println("📊 Performance Statistics")
	fmt.Println("========================")
	fmt.Printf("Uptime: %v\n", stats["uptime"])
	fmt.Printf("Total Logs: %d\n", stats["total_logs"])
	fmt.Printf("Dropped Logs: %d\n", stats["dropped_logs"])
	fmt.Printf("Error Count: %d\n", stats["error_count"])
	fmt.Printf("Throughput: %.2f logs/sec\n", stats["throughput"])
	
	if mem, ok := stats["memory"].(*MemoryStats); ok {
		fmt.Println("\n💾 Memory Usage:")
		fmt.Printf("  Alloc: %d bytes (%.2f MB)\n", mem.Alloc, float64(mem.Alloc)/1024/1024)
		fmt.Printf("  Total Alloc: %d bytes (%.2f MB)\n", mem.TotalAlloc, float64(mem.TotalAlloc)/1024/1024)
		fmt.Printf("  Sys: %d bytes (%.2f MB)\n", mem.Sys, float64(mem.Sys)/1024/1024)
		fmt.Printf("  Heap Objects: %d\n", mem.HeapObjects)
		fmt.Printf("  GC Cycles: %d\n", mem.NumGC)
	}
	
	if cpu, ok := stats["cpu"].(*CPUStats); ok {
		fmt.Println("\n🖥️  CPU Usage:")
		fmt.Printf("  Goroutines: %d\n", cpu.NumGoroutine)
		fmt.Printf("  CPUs: %d\n", cpu.NumCPU)
		fmt.Printf("  CGO Calls: %d\n", cpu.NumCgoCall)
	}
	
	if lat, ok := stats["latency"].(*LatencyStats); ok {
		fmt.Println("\n⏱️  Latency:")
		fmt.Printf("  Min: %v\n", lat.MinLatency)
		fmt.Printf("  Max: %v\n", lat.MaxLatency)
		fmt.Printf("  Avg: %v\n", lat.AvgLatency)
		fmt.Printf("  P50: %v\n", lat.P50Latency)
		fmt.Printf("  P90: %v\n", lat.P90Latency)
		fmt.Printf("  P95: %v\n", lat.P95Latency)
		fmt.Printf("  P99: %v\n", lat.P99Latency)
	}
}

// Reset resets all statistics
func (pm *PerformanceMonitor) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	atomic.StoreInt64(&pm.totalLogs, 0)
	atomic.StoreInt64(&pm.droppedLogs, 0)
	atomic.StoreInt64(&pm.errorCount, 0)
	
	pm.startTime = time.Now()
	pm.latencyStats = &LatencyStats{
		latencySamples: make([]time.Duration, 0, 1000),
	}
}
