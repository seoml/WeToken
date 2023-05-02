package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Kdb *gorm.DB
)

type TextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}
type uKey struct {
	Id    int    `json:"id"`
	Uid   string `json:"uid"`
	Token string `json:"token"`
	State int    `json:"state"`
	Time  string `json:"time"`
}

func ChatToken(c *gin.Context) {
	tokenFile, _ := os.OpenFile("./config/Token.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer tokenFile.Close()
	var msg TextMessage
	if err := c.ShouldBindXML(&msg); err != nil {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	resp := TextMessage{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   msg.CreateTime,
		MsgType:      "text",
		Content:      "咕咕咕...",
	}
	mc := msg.Content

	switch {
	case strings.Contains(mc, "收到不支持的消息类型"):
		resp.Content = "您发送的消息暂时无法处理~可以尝试发送给客服"
		c.XML(http.StatusOK, resp)
		return
	case strings.Contains(mc, "客服"):
		kf := viper.GetString("kf")
		resp.Content = kf
		c.XML(http.StatusOK, resp)
		return
	}

	if strings.ToLower(msg.MsgType) == "text" {
		formUser := msg.FromUserName
		var tokenNum int64
		var ut uKey
		var tokenY string

		Kdb.Table("tokens").Where("uid = ?", formUser).Find(&ut).Count(&tokenNum)

		if tokenNum >= 1 { //判断有token则返回结果
			tokenY = ut.Token
			resp.Content = "你的Token为：" + tokenY
			c.XML(http.StatusOK, resp)
			return
		} else { //否则为其生成一个新的token
			ut.Uid = msg.FromUserName
			ut.State = 0
			currentTime := time.Now()
			formattedTime := currentTime.Format("2006-01-02 15:04:05")
			ut.Time = formattedTime                               //记录生成token的时间和日期
			milliTimestamp := time.Now().UnixNano() / 1e6         //1、当前时间戳
			data := fmt.Sprintf("%s%d", formUser, milliTimestamp) //2、用户id+时间戳
			hash := sha256.Sum256([]byte(data))                   //3、生成哈希
			ut.Token = hex.EncodeToString(hash[:])                //4、设置token。（也可以自己把这段改良一下）
			Kdb.Table("tokens").Create(&ut)
			resp.Content = "你的Token为：" + ut.Token
			tokenFile.WriteString("生成Token：" + ut.Token + "丨来源用户：" + formUser + "\r\n")
			c.XML(http.StatusOK, resp)
			return
		}
	}
	hello := viper.GetString("hello")
	resp.Content = hello
	c.XML(http.StatusOK, resp)
	return
}
