package model

const (
	// X_TOKEN 配置中http部分 增加underscores_in_headers on; 配置
	// 用减号-替代下划线符号_，避免这种变态问题。nginx默认忽略掉下划线可能有些原因。
	X_TOKEN   = "X-TOKEN"
	CURR_USER = "CURR_USER"
)
