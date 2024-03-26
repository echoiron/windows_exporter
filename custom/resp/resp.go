package resp

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus-community/windows_exporter/custom/util"
	"strings"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (r Resp) WithMsg(msg string) Resp {
	r.Msg = msg
	if r.Data == nil {
		r.Data = msg
	}
	return r
}

func (r Resp) WithMsgF(msg string, a ...interface{}) Resp {
	r.Msg = fmt.Sprintf(msg, a)
	return r
}

func (r Resp) WithData(data interface{}) Resp {
	r.Data = data
	return r
}

func (r Resp) Error() string {
	return fmt.Sprintf("code:%d Message:%s data:%v", r.Code, r.Msg, r.Data)
}

// TranslateError 翻译validate验证错误
func (r Resp) TranslateError(err error) Resp {
	translatedErrors := make([]string, 0)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			translatedErrors = append(translatedErrors, validationError.Translate(util.Trans))
		}
	} else {
		translatedErrors = append(translatedErrors, err.Error())
	}
	r.Msg = strings.Join(translatedErrors, ",")
	return r
}

func (r Resp) WithError(err error) Resp {
	if err != nil {
		r.Data = err.Error()
	}
	return r
}

func (r Resp) ErrorDetail(err error) Resp {
	if err != nil {
		r.Data = err.Error()
	}
	return r
}

func (r Resp) String() string {
	bs, _ := json.Marshal(r)
	result := string(bs)
	return result
}

func (r Resp) Bytes() []byte {
	bs, _ := json.Marshal(r)
	return bs
}
