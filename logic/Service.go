package logic

import (
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

	return r.Run(":" + strconv.Itoa(G_config.ServicePort))
}

// 全量推送POST msg={}
func pushAll(c *gin.Context) {
	items := c.PostForm("items")

	var msgArr []json.RawMessage
	if err := json.Unmarshal([]byte(items), &msgArr); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err := G_gateConnMgr.PushAll(msgArr)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "pushed")
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

	err := G_gateConnMgr.PushRoom(room, msgArr)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "pushed")
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
