package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DecodeErr(err error, data interface{}) (int, string, interface{}) {
	if err == nil {
		return OK.Code, OK.Msg, data
	}

	switch typed := err.(type) {
	case Resp:
		if typed.Code == 200 {
			return OK.Code, OK.Msg, data
		}
		return typed.Code, typed.Msg, typed.Data
	case *Resp:
		return typed.Code, typed.Msg, typed.Data
	}
	return Unauthorized.Code, Unauthorized.Msg, err.Error()
}

// 发送json响应信息
func SendJsonResponse(c *gin.Context, err error, data interface{}) {
	code, message, decodeErr := DecodeErr(err, data)
	resp := Resp{
		Code: code,
		Msg:  message,
		Data: decodeErr,
	}
	if data != nil {
		resp.Data = data
	}
	c.JSON(http.StatusOK, resp)
	return
}
