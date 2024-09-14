package constant

type ContextKeyType string

// cache keys
const (
	OrderCacheKey = "order"
)

const (
	TimeFormat                  = "2006-01-02 15:04:05"
	UserSession  ContextKeyType = "user-session"
	DefaultLimit                = 30
)
