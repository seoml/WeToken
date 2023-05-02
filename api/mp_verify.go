package api

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sort"
	"strings"
)

var wcToken string

func MpVerify(c *gin.Context) {
	wcToken = viper.GetString("wechat.token")
	// 获取微信传递的参数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 将三个参数按字典序排序
	params := []string{wcToken, nonce, timestamp}
	sort.Strings(params)

	// 将三个参数拼接成一个字符串
	str := strings.Join(params, "")

	// 对拼接后的字符串进行sha1加密
	sha := sha1.New()
	sha.Write([]byte(str))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// 将加密后的字符串与微信传递的signature参数进行比较
	if hash == signature {
		// 如果一致，返回echostr参数
		c.String(200, echostr)
	} else {
		// 如果不一致，返回错误信息
		c.String(401, "Token效验失败")
		fmt.Println("Token效验失败了~")
	}
}
