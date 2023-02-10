package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RespCode int

const (
	RespCodeOKCode               RespCode = 0
	RespCodeOKMsg                string   = "ok"
	RespCodeERCode               RespCode = 1
	RespCodeERMsg                string   = "failed to respond"
	RespCodeInvalidInputCode     RespCode = 2
	RespCodeInvalidInputMsg      string   = "invalid input"
	RespCodeProcessingFailedCode RespCode = 2
	RespCodeProcessingFailedMsg  string   = "processing failed"
)

type RespAPIModel struct {
	Code RespCode `json:"code"`
	Msg  string   `json:"msg"`
}

func (r *RespAPIModel) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, RespCodeERCode, RespCodeERMsg)
	}
	return string(data)
}

func NewResp(code RespCode, msg string) *RespAPIModel {
	return &RespAPIModel{
		Code: code,
		Msg:  msg,
	}
}

func RespAPIFailed(w http.ResponseWriter) {
	_, _ = w.Write([]byte(NewResp(RespCodeERCode, RespCodeERMsg).String()))
}

func RespAPIOk(w http.ResponseWriter) {
	_, _ = w.Write([]byte(NewResp(RespCodeOKCode, RespCodeOKMsg).String()))
}

func RespAPIInvalidInput(w http.ResponseWriter) {
	_, _ = w.Write([]byte(NewResp(RespCodeInvalidInputCode, RespCodeInvalidInputMsg).String()))
}

func RespAPIProcessingFailed(w http.ResponseWriter) {
	_, _ = w.Write([]byte(NewResp(RespCodeProcessingFailedCode, RespCodeProcessingFailedMsg).String()))
}

func RespRedirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
