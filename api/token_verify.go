package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func TokenVerify(c *gin.Context) {
	token := c.PostForm("token")
	var num int64
	var data uKey
	Kdb.Table("tokens").Where("token = ?", token).Limit(1).Find(&data).Count(&num)
	if num == 0 {
		c.JSON(200, gin.H{"code": 404, "verify": "no", "msg": "你输入的token不存在！"})
		return

	}
	switch data.State {
	case 0:
		c.JSON(200, gin.H{"code": 200, "verify": "ok", "msg": "验证成功！", "data": data})
		return
	case 1:
		c.JSON(200, gin.H{"code": 250, "verify": "no", "msg": "你的Token已被禁用！"})
		return
	}
	c.JSON(200, gin.H{"code": 404, "verify": "no", "msg": "你输入的token不存在！"})

	fmt.Println("触发了兜底返回，可能遇到了什么问题。\r\n用户输入的Token为：", token)
	return
}
