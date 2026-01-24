package business

import (
	"qauth-server/internal/config"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/services"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	userService    *services.UserService
	counterService *services.CounterService
	cacheService   *services.CacheService
}

func NewDashboardHandler(
	userService *services.UserService,
	counterService *services.CounterService,
	cacheService *services.CacheService,
) *DashboardHandler {
	return &DashboardHandler{
		userService:    userService,
		counterService: counterService,
		cacheService:   cacheService,
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	// 获取用户总数
	userCnt, err := h.userService.UserCount()
	if err != nil {
		utilities.GetLogger().Error("failed to get total user count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}
	userCounter, err := h.counterService.GetRecentCounter(string(config.CounterTotalUser), 9, nil)
	if err != nil {
		utilities.GetLogger().Error("failed to get total user trend", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 组装响应数据
	userTrend := make([]int64, 0)
	for _, cnt := range userCounter {
		userTrend = append(userTrend, cnt.Count)
	}
	userTrend = append(userTrend, userCnt)

	// 获取最近10天的认证用户数量趋势
	authCounter, err := h.counterService.GetRecentCounter(string(config.CounterAuthUser), 9, &services.ZERO)
	if err != nil {
		utilities.GetLogger().Error("failed to get recent auth user count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}
	userAuthCnt, err := h.cacheService.GetKeyValueAsInt64("todays-authcount")
	if err != nil {
		utilities.GetLogger().Error("failed to get todays auth count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 组装响应数据
	userAuthTrend := make([]int64, 0)
	for _, cnt := range authCounter {
		userAuthTrend = append(userAuthTrend, cnt.Count)
	}
	userAuthTrend = append(userAuthTrend, userAuthCnt)

	// 获取最近 10 周的 OAuth 应用注册数量趋势
	oauthAppCounter, err := h.counterService.GetRecentCounter(string(config.CounterOAuthApp), 9, nil)
	if err != nil {
		utilities.GetLogger().Error("failed to get recent oauth app count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}
	oauthAppCnt, err := h.cacheService.GetKeyValueAsInt64("weekly-oauthapp-count")
	if err != nil {
		utilities.GetLogger().Error("failed to get weekly oauth app count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 组装响应数据
	oauthAppTrend := make([]int64, 0)
	for _, cnt := range oauthAppCounter {
		oauthAppTrend = append(oauthAppTrend, cnt.Count)
	}
	oauthAppTrend = append(oauthAppTrend, oauthAppCnt)

	// 获取当前活跃 User 数量
	realtimeActiveCnt, err := h.cacheService.CountPrefixKeys("oauth2_token:atk:")
	if err != nil {
		utilities.GetLogger().Error("failed to get active user count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}
	maxActiveTrend, err := h.counterService.GetRecentCounter(string(config.CounterDailyMaxActiveUser), 9, &services.ZERO)
	if err != nil {
		utilities.GetLogger().Error("failed to get active user trend", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 组装响应数据
	realtimeActiveCnt = realtimeActiveCnt - 1 // 减去缓存中的空 key
	activeTrendData := make([]int64, 0)
	for _, cnt := range maxActiveTrend {
		activeTrendData = append(activeTrendData, cnt.Count)
	}
	activeTrendData = append(activeTrendData, realtimeActiveCnt)

	response.HandlerSuccess(c, gin.H{
		"user_count":        userCnt,
		"user_trend":        userTrend,
		"user_auth_count":   userAuthCnt,
		"user_auth_trend":   userAuthTrend,
		"oauth_app_count":   oauthAppCnt,
		"oauth_app_trend":   oauthAppTrend,
		"active_user_count": realtimeActiveCnt,
		"active_user_trend": activeTrendData,
	})
}
