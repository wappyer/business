package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"

	. "business/common"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyLogWriter) Write(data []byte) (int, error) {
	b.body.Write(data)
	return b.ResponseWriter.Write(data)
}

func (b bodyLogWriter) WriteString(s string) (int, error) {
	b.body.WriteString(s)
	return b.ResponseWriter.WriteString(s)
}

func ReqLog() gin.HandlerFunc {
	return func(g *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: g.Writer}
		g.Writer = bodyLogWriter

		start := time.Now().UnixNano() / 1e6
		nowNano := IntToStr(int(time.Now().UnixNano()))
		ReqNo = time.Now().Format("20060102-150405-") + nowNano[len(nowNano)-9:len(nowNano)]

		ReqLogChan <- &ReqLogForChan{
			ReqNo:      ReqNo,
			UserId:     TokenInfo.UserId,
			Router:     g.Request.URL.Path,
			Method:     g.Request.Method,
			Agent:      g.Request.UserAgent(),
			Param:      g.Request.RequestURI,
			Ip:         g.ClientIP(),
			CreateTime: GetNow(),
		}

		g.Next()

		// 记录接口返回
		respBody := bodyLogWriter.body.String()
		resp := RespBody{}
		_ = json.Unmarshal([]byte(respBody), &resp)
		data, _ := json.Marshal(resp.Data)
		cost := float64(time.Now().UnixNano()/1e6-start) / 1000
		ReqLogChan <- &ReqLogForChan{
			ReqNo:    ReqNo,
			UserId:   TokenInfo.UserId,
			Cost:     cost,
			HttpCode: g.Writer.Status(),
			Code:     resp.Code,
			Msg:      resp.Msg,
			Data:     bytes.NewBuffer(data).String(),
		}
	}
}
