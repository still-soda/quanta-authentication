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
	oauthService   *services.OAuthService
}

func NewDashboardHandler(
	userService *services.UserService,
	counterService *services.CounterService,
	cacheService *services.CacheService,
	oauthService *services.OAuthService,
) *DashboardHandler {
	return &DashboardHandler{
		userService:    userService,
		counterService: counterService,
		cacheService:   cacheService,
		oauthService:   oauthService,
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
	oauthAppCnt, err := h.oauthService.CountClients()
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

// GetUserDistributionByRole 获取用户按角色的分布数据
func (h *DashboardHandler) GetUserDistributionByRole(c *gin.Context) {
	roleCountMap, err := h.userService.GetUserCountByRole()
	if err != nil {
		utilities.GetLogger().Error("failed to get user distribution by role", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 定义颜色映射
	colorMap := map[string]string{
		"系统超级管理员": "#f97316", // orange
		"系统管理员":   "#3b82f6", // blue
		"系统用户":    "#10b981", // green
		"未分配角色":   "#9ca3af", // gray
	}

	defaultColors := []string{"#8b5cf6", "#06b6d4", "#ec4899", "#84cc16", "#f59e0b"}

	labels := make([]string, 0)
	data := make([]int64, 0)
	colors := make([]string, 0)
	colorIndex := 0

	for roleName, count := range roleCountMap {
		if count > 0 {
			labels = append(labels, roleName)
			data = append(data, count)

			// 使用预定义颜色或默认颜色
			if color, ok := colorMap[roleName]; ok {
				colors = append(colors, color)
			} else {
				colors = append(colors, defaultColors[colorIndex%len(defaultColors)])
				colorIndex++
			}
		}
	}

	// 如果没有数据，返回默认值
	if len(labels) == 0 {
		labels = []string{"暂无数据"}
		data = []int64{1}
		colors = []string{"#e5e7eb"}
	}

	response.HandlerSuccess(c, gin.H{
		"labels": labels,
		"data":   data,
		"colors": colors,
	})
}

// GetAuthTrend 获取最近认证用户数量趋势
func (h *DashboardHandler) GetAuthTrend(c *gin.Context) {
	rangeType := c.DefaultQuery("range", "weekly")

	var days int
	switch rangeType {
	case "weekly":
		days = 7
	case "half-weekly":
		days = 15
	case "monthly":
		days = 30
	default:
		utilities.GetLogger().Error("invalid range type for auth trend", "range", rangeType)
		response.HandlerError(c, app_error.ErrBadRequest)
		return
	}

	authCounter, err := h.counterService.GetRecentCounter(string(config.CounterAuthUser), days-1, &services.ZERO)
	if err != nil {
		utilities.GetLogger().Error("failed to get auth user trend", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	authCountToday, err := h.cacheService.GetKeyValueAsInt64("todays-authcount")
	if err != nil {
		utilities.GetLogger().Error("failed to get today's auth user count", "error", err.Error())
		response.HandlerError(c, app_error.ErrInternalServerError)
		return
	}

	// 组装响应数据
	authTrend := make([]int64, 0)
	for _, cnt := range authCounter {
		authTrend = append(authTrend, cnt.Count)
	}
	authTrend = append(authTrend, authCountToday)

	response.HandlerSuccess(c, gin.H{
		"auth_trend": authTrend,
	})
}
