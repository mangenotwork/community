package handler

import (
	"community/internal/local/datasvc"
	"community/internal/local/registersvc"
	"github.com/gin-gonic/gin"
)

func Handler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": 0,
	})
}

// 测试add
func AddT(ctx *gin.Context) {
	name := "test add "
	b, _ := datasvc.Add(name, name, []byte(name))
	registersvc.NodeTable.BroadcastData(b)
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}
