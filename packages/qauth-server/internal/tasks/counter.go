package tasks

import (
	"qauth-server/internal/config"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"

	"github.com/robfig/cron/v3"
)

type CounterTask struct {
	cacheService   *services.CacheService
	counterService *services.CounterService
	userService    *services.UserService
}

func NewCounterTask(
	counterService *services.CounterService,
	cacheService *services.CacheService,
	userService *services.UserService,
) *CounterTask {
	return &CounterTask{
		counterService: counterService,
		cacheService:   cacheService,
		userService:    userService,
	}
}

// Register 注册定时任务，返回注销函数
func (t *CounterTask) Register(cronScheduler *cron.Cron) (func(), error) {
	authUserCountTaskID, err := cronScheduler.AddFunc("0 0 0 * * *", t.SaveAuthUserCount)
	if err != nil {
		return nil, err
	}

	weeklyOAuthAppCountTaskID, err := cronScheduler.AddFunc("0 0 0 * * 0", t.SaveWeeklyOAuthAppCount)
	if err != nil {
		return nil, err
	}

	dailyMaxActiveUserCountTaskID, err := cronScheduler.AddFunc("0 0 0 * * *", t.SaveDailyMaxActiveUserCount)
	if err != nil {
		return nil, err
	}

	return func() {
		cronScheduler.Remove(authUserCountTaskID)
		cronScheduler.Remove(weeklyOAuthAppCountTaskID)
		cronScheduler.Remove(dailyMaxActiveUserCountTaskID)
	}, nil
}

// 保存每日认证用户数量到计数器
func (t *CounterTask) SaveAuthUserCount() {
	authUserCnt, err := t.cacheService.GetKeyValueAsInt64("todays-authcount")
	if err != nil {
		utilities.GetLogger().Error("faild to save auth user count", "error", err.Error(), "cnt", authUserCnt)
		return
	}

	t.counterService.CreateCounter(string(config.CounterAuthUser), authUserCnt)
}

// 保存每周注册的 OAuth 应用数量到计数器
func (t *CounterTask) SaveWeeklyOAuthAppCount() {
	oauthAppCnt, err := t.cacheService.GetKeyValueAsInt64("weekly-oauthapp-count")
	if err != nil {
		utilities.GetLogger().Error("failed to save weekly oauth app count", "error", err.Error(), "cnt", oauthAppCnt)
		return
	}

	t.counterService.CreateCounter(string(config.CounterOAuthApp), oauthAppCnt)
}

// 保存每日最大活跃用户数量
func (t *CounterTask) SaveDailyMaxActiveUserCount() {
	activeUserCnt, err := t.cacheService.GetKeyValueAsInt64("todays-max-active-user-count")
	if err != nil {
		utilities.GetLogger().Error("failed to save daily max active user count", "error", err.Error(), "cnt", activeUserCnt)
		return
	}

	t.counterService.CreateCounter(string(config.CounterDailyMaxActiveUser), activeUserCnt)

	// 重置当天最大活跃用户数缓存
	t.cacheService.SetKeyValue("todays-max-active-user-count", "0", 24*3600)
}

// 保存用户数量到计数器
func (t *CounterTask) SaveUserCount() {
	userCnt, err := t.userService.UserCount()
	if err != nil {
		utilities.GetLogger().Error("failed to save user count", "error", err.Error(), "cnt", userCnt)
		return
	}

	t.counterService.CreateCounter(string(config.CounterTotalUser), userCnt)
}
