package cronManager

import (
	"sync"

	"github.com/robfig/cron/v3"
	"wails-router-link/service/types"
)

//type CronInterface interface {
//	// FindCronList Note Mortal 2026/4/28 15:03 获取所有的定时任务
//	FindCronList() map[string]*taskManager
//	// AddTaskByFuncWithSecond Note Mortal 2026/4/28 15:03 添加Task 方法形式以秒的形式加入
//	AddTaskByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) // 添加Task Func以秒的形式加入
//	// AddTaskByJobWithSeconds Note Mortal 2026/4/28 15:03 添加Task 接口形式以秒的形式加入
//	AddTaskByJobWithSeconds(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)
//	// AddTaskByFunc Note Mortal 2026/4/28 15:03 通过函数的方法添加任务
//	AddTaskByFunc(cronName string, spec string, task func(), taskName string, option ...cron.Option) (cron.EntryID, error)
//	// AddTaskByJob Note Mortal 2026/4/28 15:03 通过接口的方法添加任务 要实现一个带有 Run方法的接口触发
//	AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)
//	// FindCron Note Mortal 2026/4/28 15:03 获取对应taskName的cron 可能会为空
//	FindCron(cronName string) (*taskManager, bool)
//	// StartCron Note Mortal 2026/4/28 15:03 开始任务
//	StartCron(cronName string)
//	// StopCron Note Mortal 2026/4/28 15:03 停止任务
//	StopCron(cronName string)
//	// FindTask Note Mortal 2026/4/28 15:03 查找指定cron下的指定task
//	FindTask(cronName string, taskName string) (*task, bool)
//	// RemoveTask Note Mortal 2026/4/28 15:03 根据id删除指定cron下的指定task
//	RemoveTask(cronName string, id int)
//	// RemoveTaskByName Note Mortal 2026/4/28 15:03 根据taskName删除指定cron下的指定task
//	RemoveTaskByName(cronName string, taskName string)
//	// Clear Note Mortal 2026/4/28 15:03 清除指定任务
//	Clear(cronName string)
//	// Close Note Mortal 2026/4/28 15:03 停止所有的cron
//	Close()
//}

type task struct {
	EntryID  cron.EntryID
	Spec     string
	TaskName string
}

type taskManager struct {
	corn  *cron.Cron
	tasks map[cron.EntryID]*task
}

// CronManager 定时任务管理
type CronManager struct {
	cronList map[string]*taskManager
	sync.Mutex
}

func NewTimerTask(config *types.AppConfig) *CronManager {
	config.Cron.Enable = true
	return &CronManager{cronList: make(map[string]*taskManager)}
}

// AddTaskByFunc Mortal 2026/4/28 15:07 通过函数的方法添加任务
func (t *CronManager) AddTaskByFunc(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.cronList[cronName]; !ok {
		tasks := make(map[cron.EntryID]*task)
		t.cronList[cronName] = &taskManager{
			corn:  cron.New(option...),
			tasks: tasks,
		}
	}
	id, err := t.cronList[cronName].corn.AddFunc(spec, fun)
	t.cronList[cronName].corn.Start()
	t.cronList[cronName].tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// AddTaskByFuncWithSecond Mortal 2026/4/28 15:07 添加Task 方法形式以秒的形式加入
func (t *CronManager) AddTaskByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	option = append(option, cron.WithSeconds())
	if _, ok := t.cronList[cronName]; !ok {
		tasks := make(map[cron.EntryID]*task)
		t.cronList[cronName] = &taskManager{
			corn:  cron.New(option...),
			tasks: tasks,
		}
	}
	id, err := t.cronList[cronName].corn.AddFunc(spec, fun)
	t.cronList[cronName].corn.Start()
	t.cronList[cronName].tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// AddTaskByJob Mortal 2026/4/28 15:07 通过接口的方法添加任务 要实现一个带有 Run方法的接口触发
func (t *CronManager) AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.cronList[cronName]; !ok {
		tasks := make(map[cron.EntryID]*task)
		t.cronList[cronName] = &taskManager{
			corn:  cron.New(option...),
			tasks: tasks,
		}
	}
	id, err := t.cronList[cronName].corn.AddJob(spec, job)
	t.cronList[cronName].corn.Start()
	t.cronList[cronName].tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// AddTaskByJobWithSeconds Mortal 2026/4/28 15:07  添加Task 接口形式以秒的形式加入
func (t *CronManager) AddTaskByJobWithSeconds(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	t.Lock()
	defer t.Unlock()
	option = append(option, cron.WithSeconds())
	if _, ok := t.cronList[cronName]; !ok {
		tasks := make(map[cron.EntryID]*task)
		t.cronList[cronName] = &taskManager{
			corn:  cron.New(option...),
			tasks: tasks,
		}
	}
	id, err := t.cronList[cronName].corn.AddJob(spec, job)
	t.cronList[cronName].corn.Start()
	t.cronList[cronName].tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// FindCron Mortal 2026/4/28 15:07 获取对应taskName的cron 可能会为空
func (t *CronManager) FindCron(cronName string) (*taskManager, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.cronList[cronName]
	return v, ok
}

// FindTask Mortal 2026/4/28 15:07 查找指定cron下的指定task
func (t *CronManager) FindTask(cronName string, taskName string) (*task, bool) {
	t.Lock()
	defer t.Unlock()
	v, ok := t.cronList[cronName]
	if !ok {
		return nil, ok
	}
	for _, t2 := range v.tasks {
		if t2.TaskName == taskName {
			return t2, true
		}
	}
	return nil, false
}

// FindCronList Mortal 2026/4/28 15:07 获取所有的任务列表
func (t *CronManager) FindCronList() map[string]*taskManager {
	t.Lock()
	defer t.Unlock()
	return t.cronList
}

// StartCron 开始任务
// StartCron Mortal 2026/4/28 15:07 开始任务
func (t *CronManager) StartCron(cronName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Start()
	}
}

// StopCron Mortal 2026/4/28 15:07 停止任务
func (t *CronManager) StopCron(cronName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Stop()
	}
}

// RemoveTask Mortal 2026/4/28 15:07 根据id删除指定cron下的指定task
func (t *CronManager) RemoveTask(cronName string, id int) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Remove(cron.EntryID(id))
		delete(v.tasks, cron.EntryID(id))
	}
}

// RemoveTaskByName Mortal 2026/4/28 15:07 根据taskName删除指定cron下的指定task
func (t *CronManager) RemoveTaskByName(cronName string, taskName string) {
	fTask, ok := t.FindTask(cronName, taskName)
	if !ok {
		return
	}
	t.RemoveTask(cronName, int(fTask.EntryID))
}

// Clear Mortal 2026/4/28 15:07 清除指定任务
func (t *CronManager) Clear(cronName string) {
	t.Lock()
	defer t.Unlock()
	if v, ok := t.cronList[cronName]; ok {
		v.corn.Stop()
		delete(t.cronList, cronName)
	}
}

// Close Mortal 2026/4/28 15:07 停止所有定时任务
func (t *CronManager) Close() {
	t.Lock()
	defer t.Unlock()
	for _, v := range t.cronList {
		v.corn.Stop()
	}
}
