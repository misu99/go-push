package gateway

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitService() error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/push/all", pushAll)
	r.POST("/push/room", pushRoom)
	r.GET("/stats", stats)

	// TLS证书解析验证
	if _, err := tls.LoadX509KeyPair(G_config.ServerPem, G_config.ServerKey); err != nil {
		//return common.ERR_CERT_INVALID
		return err
	}

	return r.RunTLS(":" + strconv.Itoa(G_config.ServicePort), G_config.ServerPem, G_config.ServerKey)
}

// 全量推送POST msg={}
func pushAll(c *gin.Context) {
	items := c.PostForm("items")

	var msgArr []json.RawMessage
	if err := json.Unmarshal([]byte(items), &msgArr); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for msgIdx := range msgArr {
		_ = G_merger.PushAll(&msgArr[msgIdx])
	}
}

// 房间推送POST room=xxx&msg
func pushRoom(c *gin.Context) {
	room := c.PostForm("room")
	items := c.PostForm("items")

	var msgArr []json.RawMessage
	if err := json.Unmarshal([]byte(items), &msgArr); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for msgIdx := range msgArr {
		_ = G_merger.PushRoom(room, &msgArr[msgIdx])
	}
}

// 统计查询
func stats(c *gin.Context) {
	data, err := G_stats.Dump()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, string(data))
}
