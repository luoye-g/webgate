package ctx

import "github.com/gin-gonic/gin"

const (
	UserIDKey = "user_id"
)

func GetUserID(c *gin.Context) int {
	return c.GetInt(UserIDKey)
}

func SetUserID(c *gin.Context, userID uint64) {
	c.Set(UserIDKey, userID)
}
