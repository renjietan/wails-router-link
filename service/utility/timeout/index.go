package utility

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// timerInfo 存储定时器详细信息
type timerInfo struct {
	name      string
	duration  time.Duration
	startTime time.Time
	endTime   time.Time
	timer     *time.Timer
	fn        func()
}

// timerCommand 定义管理器支持的命令类型
type timerCommand struct {
	cmdType  string        // 命令类型
	name     string        // 定时器名称
	duration time.Duration // 定时器时长
	fn       func()        // 回调函数
	done     chan struct{} // 操作完成通知通道
	result   chan any      // 结果返回通道
}

// ResetOption 重置选项类型
type ResetOption func(*resetOptions)

type resetOptions struct {
	duration time.Duration
	fn       func()
}

// TimeoutManager 使用通道模式实现的并发安全定时器管理器
type TimeoutManager struct {
	cmdChan   chan timerCommand // 命令通道
	stopped   atomic.Bool       // 停止标志（atomic保证并发安全）
	wg        sync.WaitGroup    // 等待组
	closeOnce sync.Once         // 确保只关闭一次
}

// NewTimeoutManager 创建新的定时器管理器
func NewTimeoutManager() *TimeoutManager {
	tm := &TimeoutManager{
		cmdChan: make(chan timerCommand, 1024), // 缓冲通道，避免阻塞
	}

	tm.wg.Add(1)
	go tm.managerLoop()

	return tm
}

// ==================== 核心管理循环 ====================

func (tm *TimeoutManager) managerLoop() {
	defer tm.wg.Done()

	timers := make(map[string]*timerInfo)
	stats := struct {
		totalCreated int64
		totalExpired int64
		totalStopped int64
		totalReset   int64
	}{}

	for cmd := range tm.cmdChan {
		switch cmd.cmdType {

		case "set":
			// 停止同名定时器
			if oldInfo, exists := timers[cmd.name]; exists {
				if oldInfo.timer.Stop() {
					stats.totalStopped++
				}
				delete(timers, cmd.name)
			}

			// 创建新的定时器信息
			info := &timerInfo{
				name:      cmd.name,
				duration:  cmd.duration,
				startTime: time.Now(),
				endTime:   time.Now().Add(cmd.duration),
				fn:        cmd.fn,
			}

			// 创建定时器
			info.timer = time.AfterFunc(cmd.duration, func() {
				// 执行用户回调（捕获panic）
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("定时器 %s 回调函数触发 panic: %v\n", cmd.name, r)
					}
				}()

				if cmd.fn != nil {
					cmd.fn()
				}

				// 定时器到期后，通知管理器清理
				select {
				case tm.cmdChan <- timerCommand{
					cmdType: "expired",
					name:    cmd.name,
				}:
				default:
					// 非阻塞发送，避免死锁
				}
			})

			timers[cmd.name] = info
			stats.totalCreated++

			if cmd.done != nil {
				close(cmd.done)
			}

		case "reset":
			// 重置定时器
			if oldInfo, exists := timers[cmd.name]; exists {
				// 停止旧定时器
				if oldInfo.timer.Stop() {
					stats.totalStopped++
				}

				// 使用新函数或原函数
				newFn := cmd.fn
				newDuration := cmd.duration

				if newFn == nil {
					newFn = oldInfo.fn // 使用原函数
				}
				if newDuration == 0 {
					newDuration = oldInfo.duration // 使用原时长
				}

				// 创建新定时器信息
				info := &timerInfo{
					name:      cmd.name,
					duration:  newDuration,
					startTime: time.Now(),
					endTime:   time.Now().Add(newDuration),
					fn:        newFn,
				}

				// 创建新定时器
				info.timer = time.AfterFunc(newDuration, func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("定时器 %s 回调函数触发 panic: %v\n", cmd.name, r)
						}
					}()

					if newFn != nil {
						newFn()
					}

					// 定时器到期后，通知管理器清理
					select {
					case tm.cmdChan <- timerCommand{
						cmdType: "expired",
						name:    cmd.name,
					}:
					default:
					}
				})

				timers[cmd.name] = info
				stats.totalReset++

				if cmd.result != nil {
					cmd.result <- struct {
						OldDuration  time.Duration
						NewDuration  time.Duration
						UsedOriginal bool
					}{
						OldDuration:  oldInfo.duration,
						NewDuration:  newDuration,
						UsedOriginal: cmd.fn == nil,
					}
				}
			} else {
				// 定时器不存在，返回错误
				if cmd.result != nil {
					cmd.result <- fmt.Errorf("延时器 %s 不存在", cmd.name)
				}
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "stop":
			if info, exists := timers[cmd.name]; exists {
				if info.timer.Stop() {
					stats.totalStopped++
				}
				delete(timers, cmd.name)
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "stop_all":
			for name, info := range timers {
				if info.timer.Stop() {
					stats.totalStopped++
				}
				delete(timers, name)
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "expired":
			// 定时器自然到期，清理资源并更新统计
			if _, exists := timers[cmd.name]; exists {
				delete(timers, cmd.name)
				stats.totalExpired++
			}

		case "exists":
			// 检查定时器是否存在
			_, exists := timers[cmd.name]
			if cmd.result != nil {
				cmd.result <- exists
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "get":
			// 获取定时器详细信息
			var result any
			if info, exists := timers[cmd.name]; exists {
				result = struct {
					Name      string
					Duration  time.Duration
					StartTime time.Time
					EndTime   time.Time
					Remaining time.Duration
					IsActive  bool
				}{
					Name:      info.name,
					Duration:  info.duration,
					StartTime: info.startTime,
					EndTime:   info.endTime,
					Remaining: time.Until(info.endTime),
					IsActive:  true,
				}
			} else {
				result = nil
			}

			if cmd.result != nil {
				cmd.result <- result
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "list":
			// 获取所有活跃定时器名称列表
			names := make([]string, 0, len(timers))
			for name := range timers {
				names = append(names, name)
			}

			if cmd.result != nil {
				cmd.result <- names
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "list_all":
			// 获取所有定时器详细信息
			infos := make([]struct {
				Name      string
				Duration  time.Duration
				StartTime time.Time
				EndTime   time.Time
				Remaining time.Duration
			}, 0, len(timers))

			for _, info := range timers {
				infos = append(infos, struct {
					Name      string
					Duration  time.Duration
					StartTime time.Time
					EndTime   time.Time
					Remaining time.Duration
				}{
					Name:      info.name,
					Duration:  info.duration,
					StartTime: info.startTime,
					EndTime:   info.endTime,
					Remaining: time.Until(info.endTime),
				})
			}

			if cmd.result != nil {
				cmd.result <- infos
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "stats":
			// 返回统计信息
			if cmd.result != nil {
				activeCount := len(timers)
				cmd.result <- struct {
					Active       int
					TotalCreated int64
					TotalExpired int64
					TotalStopped int64
					TotalReset   int64
				}{
					Active:       activeCount,
					TotalCreated: stats.totalCreated,
					TotalExpired: stats.totalExpired,
					TotalStopped: stats.totalStopped,
					TotalReset:   stats.totalReset,
				}
			}

			if cmd.done != nil {
				close(cmd.done)
			}

		case "quit":
			// 清理所有定时器并退出
			for name, info := range timers {
				info.timer.Stop()
				delete(timers, name)
			}

			if cmd.done != nil {
				close(cmd.done)
			}
			return
		}
	}
}

// ==================== 公共接口：基础操作 ====================

// Set 异步设置定时器（不等待确认）
func (tm *TimeoutManager) Set(name string, duration time.Duration, fn func()) {
	if tm.stopped.Load() {
		return
	}

	select {
	case tm.cmdChan <- timerCommand{
		cmdType:  "set",
		name:     name,
		duration: duration,
		fn:       fn,
	}:
		// 成功发送
	case <-time.After(100 * time.Millisecond):
		// 通道满或超时
	}
}

// SetSync 同步设置定时器（等待操作完成）
func (tm *TimeoutManager) SetSync(name string, duration time.Duration, fn func()) error {
	if tm.stopped.Load() {
		return fmt.Errorf("超时管理器已停止")
	}

	done := make(chan struct{})
	select {
	case tm.cmdChan <- timerCommand{
		cmdType:  "set",
		name:     name,
		duration: duration,
		fn:       fn,
		done:     done,
	}:
		select {
		case <-done:
			return nil
		case <-time.After(5 * time.Second):
			return fmt.Errorf("等待定时器设置确认超时")
		}
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("超时：超时管理器通道繁忙")
	}
}

// StopByName 停止指定定时器
func (tm *TimeoutManager) StopByName(name string) bool {
	if tm.stopped.Load() {
		return false
	}

	done := make(chan struct{})
	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "stop",
		name:    name,
		done:    done,
	}:
		select {
		case <-done:
			return true
		case <-time.After(1 * time.Second):
			return false
		}
	default:
		return false
	}
}

// StopAll 停止所有定时器
func (tm *TimeoutManager) StopAll() {
	if tm.stopped.Load() {
		return
	}

	done := make(chan struct{})
	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "stop_all",
		done:    done,
	}:
		<-done
	default:
		// 非阻塞操作
	}
}

// ==================== 公共接口：查询操作 ====================

// Exists 检查指定定时器是否存在
func (tm *TimeoutManager) Exists(name string) bool {
	if tm.stopped.Load() {
		return false
	}

	resultChan := make(chan any, 1)
	done := make(chan struct{})

	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "exists",
		name:    name,
		result:  resultChan,
		done:    done,
	}:
		select {
		case <-done:
			if result, ok := <-resultChan; ok {
				if exists, ok := result.(bool); ok {
					return exists
				}
			}
		case <-time.After(1 * time.Second):
		}
	default:
	}

	return false
}

// GetTimerInfo 获取指定定时器的详细信息
func (tm *TimeoutManager) GetTimerInfo(name string) (map[string]any, error) {
	if tm.stopped.Load() {
		return nil, fmt.Errorf("超时管理器已停止")
	}

	resultChan := make(chan any, 1)
	done := make(chan struct{})

	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "get",
		name:    name,
		result:  resultChan,
		done:    done,
	}:
		select {
		case <-done:
			if result, ok := <-resultChan; ok {
				if result == nil {
					return nil, fmt.Errorf("延时器 %s 不存在", name)
				}

				if info, ok := result.(struct {
					Name      string
					Duration  time.Duration
					StartTime time.Time
					EndTime   time.Time
					Remaining time.Duration
					IsActive  bool
				}); ok {
					return map[string]any{
						"name":       info.Name,
						"duration":   info.Duration,
						"start_time": info.StartTime,
						"end_time":   info.EndTime,
						"remaining":  info.Remaining,
						"is_active":  info.IsActive,
					}, nil
				}
			}
		case <-time.After(1 * time.Second):
			return nil, fmt.Errorf("获取延时器 %s 信息超时", name)
		}
	case <-time.After(100 * time.Millisecond):
		return nil, fmt.Errorf("超时：超时管理器通道繁忙")
	}

	return nil, fmt.Errorf("获取延时器 %s 信息失败", name)
}

// ListActiveTimers 获取所有活跃定时器的名称列表
func (tm *TimeoutManager) ListActiveTimers() []string {
	if tm.stopped.Load() {
		return []string{}
	}

	resultChan := make(chan any, 1)
	done := make(chan struct{})

	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "list",
		result:  resultChan,
		done:    done,
	}:
		select {
		case <-done:
			if result, ok := <-resultChan; ok {
				if names, ok := result.([]string); ok {
					return names
				}
			}
		case <-time.After(1 * time.Second):
		}
	default:
	}

	return []string{}
}

// ListAllTimersInfo 获取所有定时器的详细信息
func (tm *TimeoutManager) ListAllTimersInfo() []map[string]any {
	if tm.stopped.Load() {
		return []map[string]any{}
	}

	resultChan := make(chan any, 1)
	done := make(chan struct{})

	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "list_all",
		result:  resultChan,
		done:    done,
	}:
		select {
		case <-done:
			if result, ok := <-resultChan; ok {
				if infos, ok := result.([]struct {
					Name      string
					Duration  time.Duration
					StartTime time.Time
					EndTime   time.Time
					Remaining time.Duration
				}); ok {
					resultList := make([]map[string]any, len(infos))
					for i, info := range infos {
						resultList[i] = map[string]any{
							"name":       info.Name,
							"duration":   info.Duration,
							"start_time": info.StartTime,
							"end_time":   info.EndTime,
							"remaining":  info.Remaining,
						}
					}
					return resultList
				}
			}
		case <-time.After(1 * time.Second):
		}
	default:
	}

	return []map[string]any{}
}

// GetRemainingTime 获取定时器剩余时间
func (tm *TimeoutManager) GetRemainingTime(name string) (time.Duration, error) {
	info, err := tm.GetTimerInfo(name)
	if err != nil {
		return 0, err
	}

	if remaining, ok := info["remaining"].(time.Duration); ok {
		return remaining, nil
	}

	return 0, fmt.Errorf("获取定时器 %s 剩余时间失败", name)
}

// ==================== 公共接口：统计操作 ====================

// GetStats 获取管理器统计信息
func (tm *TimeoutManager) GetStats() (active int, created, expired, stopped, reset int64) {
	if tm.stopped.Load() {
		return 0, 0, 0, 0, 0
	}

	resultChan := make(chan any, 1)
	done := make(chan struct{})

	select {
	case tm.cmdChan <- timerCommand{
		cmdType: "stats",
		result:  resultChan,
		done:    done,
	}:
		select {
		case <-done:
			if result, ok := <-resultChan; ok {
				if stats, ok := result.(struct {
					Active       int
					TotalCreated int64
					TotalExpired int64
					TotalStopped int64
					TotalReset   int64
				}); ok {
					return stats.Active, stats.TotalCreated, stats.TotalExpired, stats.TotalStopped, stats.TotalReset
				}
			}
		case <-time.After(1 * time.Second):
		}
	default:
	}

	return -1, -1, -1, -1, -1
}

// ==================== 公共接口：重置操作 ====================

// ResetTimer 重置指定定时器（可选择性更新参数）
func (tm *TimeoutManager) ResetTimer(name string, newDuration time.Duration, newFn func()) error {
	if tm.stopped.Load() {
		return fmt.Errorf("超时管理器已停止")
	}

	resultChan := make(chan any, 1)
	done := make(chan struct{})

	select {
	case tm.cmdChan <- timerCommand{
		cmdType:  "reset",
		name:     name,
		duration: newDuration,
		fn:       newFn,
		result:   resultChan,
		done:     done,
	}:
		select {
		case <-done:
			if result, ok := <-resultChan; ok {
				if err, isErr := result.(error); isErr {
					return err
				}
				return nil
			}
		case <-time.After(2 * time.Second):
			return fmt.Errorf("重置定时器 %s 超时", name)
		}
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("超时：重置定时器 %s 通道繁忙", name)
	}

	return fmt.Errorf("重置 %s 失败", name)
}

// ResetTimerExt 重置定时器的扩展版本，允许部分参数不更新
func (tm *TimeoutManager) ResetTimerExt(name string, options ...ResetOption) error {
	if tm.stopped.Load() {
		return fmt.Errorf("超时管理器已停止")
	}

	// 解析选项
	opts := &resetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	return tm.ResetTimer(name, opts.duration, opts.fn)
}

// WithDuration 设置新的持续时间
func WithDuration(d time.Duration) ResetOption {
	return func(o *resetOptions) {
		o.duration = d
	}
}

// WithFunc 设置新的回调函数
func WithFunc(fn func()) ResetOption {
	return func(o *resetOptions) {
		o.fn = fn
	}
}

// ==================== 公共接口：管理操作 ====================

// IsActive 检查指定定时器是否活跃（Exists的别名）
func (tm *TimeoutManager) IsActive(name string) bool {
	return tm.Exists(name)
}

// IsClosed 检查管理器是否已关闭
func (tm *TimeoutManager) IsClosed() bool {
	return tm.stopped.Load()
}

// Close 优雅关闭管理器
func (tm *TimeoutManager) Close() {
	tm.closeOnce.Do(func() {
		// 设置停止标志，拒绝新请求
		tm.stopped.Store(true)

		// 发送退出命令
		done := make(chan struct{})
		select {
		case tm.cmdChan <- timerCommand{
			cmdType: "quit",
			done:    done,
		}:
			select {
			case <-done:
				// 管理器已处理退出
			case <-time.After(1 * time.Second):
				// 超时，强制继续
			}
		case <-time.After(100 * time.Millisecond):
			// 通道可能已满
		}

		// 关闭命令通道
		close(tm.cmdChan)

		// 等待管理器goroutine退出
		tm.wg.Wait()
	})
}
