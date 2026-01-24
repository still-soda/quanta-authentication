package tasks

import (
	"fmt"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"

	"github.com/robfig/cron/v3"
)

type UserTask struct {
	cacheService *services.CacheService
}

func NewUserTask(
	cacheService *services.CacheService,
) *UserTask {
	return &UserTask{
		cacheService: cacheService,
	}
}

// Register 注册定时任务，返回注销函数
func (t *UserTask) Register(cronScheduler *cron.Cron) (func(), error) {
	maxActiveUserCountTaskID, err := cronScheduler.AddFunc("0 * * * * *", t.SaveDailyMaxActiveUserCount)
	if err != nil {
		return nil, err
	}

	return func() {
		cronScheduler.Remove(maxActiveUserCountTaskID)
	}, nil
}

// 更新最大活跃用户数到缓存
func (t *UserTask) SaveDailyMaxActiveUserCount() {
	activeCnt, err := t.cacheService.CountPrefixKeys("oauth2_token:atk:")
	maxActiveCnt, err := t.cacheService.GetKeyValueAsInt64("daily-max-active-user-count")
	if err != nil {
		utilities.GetLogger().Error("failed to save daily max active user count", "error", err.Error(), "cnt", activeCnt)
		return
	}

	if activeCnt <= maxActiveCnt {
		return
	}

	// 保存新的最大活跃用户数，有效期24小时
	err = t.cacheService.SetKeyValue("daily-max-active-user-count", fmt.Sprint(activeCnt), 24*3600)
}
