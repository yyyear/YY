package YY

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SimpleGoroutineDetector 简单的 Goroutine 检测工具
type SimpleGoroutineDetector struct {
	mu           sync.RWMutex
	knownIDs     map[uint64]bool
	mainID       uint64
	creationTime map[uint64]time.Time
}

var detector *SimpleGoroutineDetector

func init() {
	detector = &SimpleGoroutineDetector{
		knownIDs:     make(map[uint64]bool),
		creationTime: make(map[uint64]time.Time),
		mainID:       getCurrentGoroutineID(),
	}
	detector.markAsKnown("main")
}

// 获取当前 goroutine ID - 高性能版本
func getCurrentGoroutineID() uint64 {
	buf := make([]byte, 32)
	buf = buf[:runtime.Stack(buf, false)]
	
	// 快速解析 "goroutine 123 [running]:"
	s := string(buf)
	if start := strings.Index(s, "goroutine "); start >= 0 {
		start += len("goroutine ")
		if end := strings.Index(s[start:], " "); end >= 0 {
			if id, err := strconv.ParseUint(s[start:start+end], 10, 64); err == nil {
				return id
			}
		}
	}
	return 0
}

// 标记当前 goroutine 为已知
func (d *SimpleGoroutineDetector) markAsKnown(purpose string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	currentID := getCurrentGoroutineID()
	if !d.knownIDs[currentID] {
		d.knownIDs[currentID] = true
		d.creationTime[currentID] = time.Now()
		fmt.Printf("🆕 新 Goroutine [ID: %d] 用途: %s\n", currentID, purpose)
	}
}

// 检查当前 goroutine 是否是新的
func (d *SimpleGoroutineDetector) isNew() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	currentID := getCurrentGoroutineID()
	return !d.knownIDs[currentID]
}

// 检查是否是主 goroutine
func (d *SimpleGoroutineDetector) isMain() bool {
	return getCurrentGoroutineID() == d.mainID
}

// 获取 goroutine 信息
func (d *SimpleGoroutineDetector) getInfo() (id uint64, isNew, isMain bool, age time.Duration) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	currentID := getCurrentGoroutineID()
	isNew = !d.knownIDs[currentID]
	isMain = currentID == d.mainID
	
	if createdAt, exists := d.creationTime[currentID]; exists {
		age = time.Since(createdAt)
	}
	
	return currentID, isNew, isMain, age
}

// QuickGoroutineCheck 快速检测函数 - 供外部调用
func QuickGoroutineCheck() string {
	id, isNew, isMain, age := detector.getInfo()
	
	status := "已知"
	if isNew {
		status = "新建"
		detector.markAsKnown("auto-detected")
	}
	
	goroutineType := "工作协程"
	if isMain {
		goroutineType = "主协程"
	}
	
	return fmt.Sprintf("[ID:%d|%s|%s|存在:%v]", id, status, goroutineType, age.Truncate(time.Millisecond))
}

// PrintGoroutineStats 打印系统统计
func PrintGoroutineStats() {
	detector.mu.RLock()
	defer detector.mu.RUnlock()
	
	fmt.Printf("\n📊 Goroutine 统计:\n")
	fmt.Printf("   系统总数: %d\n", runtime.NumGoroutine())
	fmt.Printf("   已跟踪数: %d\n", len(detector.knownIDs))
	fmt.Printf("   主协程ID: %d\n", detector.mainID)
	fmt.Printf("   当前协程: %s\n", QuickGoroutineCheck())
}
