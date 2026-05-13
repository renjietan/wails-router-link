package cronManager

import (
	"fmt"
	"testing"
	"time"
	"wails-router-link/service/types"

	"github.com/stretchr/testify/assert"
)

var job = mockJob{}

type mockJob struct{}

func (job mockJob) Run() {
	mockFunc()
}

func mockFunc() {
	time.Sleep(time.Second)
	fmt.Println("1s...")
}

func TestNewTimerTask(t *testing.T) {
	tm := NewTimerTask(&types.AppConfig{})
	_tm := tm
	{
		_, err := tm.AddTaskByFunc("1s", "@every 1s", mockFunc, "测试")
		assert.Nil(t, err, "定时任务测试失败")
		_, ok := _tm.cronList["1s"]
		assert.True(t, ok, "未找到定时任务")
	}
	{
		_, err := tm.AddTaskByJob("job", "@every 1s", job, "测试job")
		assert.Nil(t, err)
		_, ok := _tm.cronList["job"]
		assert.True(t, ok, "未找到JOB")
	}

	{
		_, ok := tm.FindCron("1s")
		if !ok {
			t.Error("未找到 1s 的任务")
		}
		_, ok = tm.FindCron("job")
		if !ok {
			t.Error("未找到 job 任务")
		}
		_, ok = tm.FindCron("none")
		if ok {
			t.Error("find none")
		}
	}
	{
		tm.Clear("1s")
		_, ok := tm.FindCron("1s")
		if ok {
			t.Error("此处不应该找到 1s ")
		}
	}
	{
		a := tm.FindCronList()
		b, c := tm.FindCron("job")
		fmt.Println(a, b, c)
	}
}
