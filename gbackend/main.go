package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/luoye-g/webgate/system"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/unrolled/secure"
	"golang.org/x/exp/rand"
)

var UserData map[string]string

func initUsers() {
	UserData = map[string]string{
		"1181895140@qq.com": "test",
	}
	go cookieDel()
}

type Session struct {
	UserInfo map[string]string
	Expire   time.Time
}

var cookieMap = make(map[string]Session)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sessionID() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]byte, 20)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func cookieDel() {
	for {
		var del []string
		for k, t := range cookieMap {

			if time.Now().Unix() > t.Expire.Unix() {
				del = append(del, k)
			}

		}
		for _, d := range del {
			delete(cookieMap, d)
		}
		time.Sleep(time.Second)
	}
}

func main() {
	system.SystemInit()

	GinHttps()
}

func GinHttps() error {

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
			return
		}
		defer ws.Close()

		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				break
			}

			// 处理接收到的消息（例如，将其广播给其他客户端）
			// ...

			// 发送响应消息
			if err := ws.WriteMessage(messageType, message); err != nil {
				break
			}
		}
	})

	r.GET("/api", func(c *gin.Context) {
		c.String(200, "test for 【%s】", "https")
	})

	r.POST("/api/login", LoginAuth)
	r.GET("/api//user_info", func(c *gin.Context) {
		// 检查cookie是否存在
		val, err := c.Cookie("user_session")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login.html")
			return
		}
		session, ok := cookieMap[val]
		if !ok {
			c.Redirect(http.StatusSeeOther, "/login.html")
			return
		}
		for k := range session.UserInfo {
			c.Set("user_name", k)
		}
	}, func(ctx *gin.Context) {

		user_name := ctx.GetString("user_name")
		ctx.JSON(200, map[string]interface{}{
			"user_name": user_name,
			"status":    200,
		})

	})
	return r.Run(":9001")
	// r.Use(TLSHandler(443))

	// return r.RunTLS(":"+strconv.Itoa(443), "./ssh/luoye-g.top.pem", "./ssh/luoye-g.top.rsa.key")
}

func TLSHandler(port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}
		c.Next()
	}
}

func CheckUserIsExist(username string) bool {
	_, isExist := UserData[username]
	return isExist
}

func CheckPassword(p1 string, p2 string) error {
	if p1 == p2 {
		return nil
	} else {
		return errors.New("password is not correct")
	}
}

func Auth(username string, password string) error {
	if isExist := CheckUserIsExist(username); isExist {
		return CheckPassword(UserData[username], password)
	} else {
		return errors.New("user is not exist")
	}
}

func LoginAuth(c *gin.Context) {
	var (
		username string
		password string
	)

	// host := c.Request.Host
	host := "luoye-g.top"
	if in, isExist := c.GetPostForm("email"); isExist && in != "" {
		username = in
	} else {
		c.HTML(http.StatusBadRequest, "/login.html", gin.H{
			"error": errors.New("必須輸入使用者名稱"),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, fmt.Sprintf("%s%s", host, "/login.html"), gin.H{
			"error": errors.New("必須輸入密碼名稱"),
		})
		return
	}
	if err := Auth(username, password); err == nil {
		expiration := time.Now().Add(time.Hour)
		c.SetCookie("user_session", "", int(expiration.Second()), "/", host, false, true)

		sessionID := sessionID()

		cookieMap[sessionID] = Session{
			UserInfo: map[string]string{username: password},
			Expire:   time.Now().Add(60 * 60 * time.Second),
		}

		c.SetCookie("user_session", sessionID, 60*60, "/", host, false, false)
		return
	} else {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s%s", host, "/login.html"))
		return
	}
}
