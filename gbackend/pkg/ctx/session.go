package ctx

import "github.com/gin-gonic/gin"

const (
	UserIDKey      = "user_id"
	UserSessionKey = "user_session"
)

func GetUserID(c *gin.Context) uint64 {
	return c.GetUint64(UserIDKey)
}

func SetUserID(c *gin.Context, userID uint64) {
	c.Set(UserIDKey, userID)
}

func SetUserSession(c *gin.Context, userSession string) {
	c.Set(UserSessionKey, userSession)
}

func GetUserSession(c *gin.Context) string {
	return c.GetString(UserSessionKey)
}
