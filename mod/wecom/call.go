package wecom

import (
	"encoding/json"
	"fmt"
	"github.com/AkvicorEdwards/glog"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// use printf to provide 'corpid' and 'corpsecret'
const getTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
const sendUrl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"

const ERCodeAccessTokenExpired = 42001

type accessTokenModel struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type tokenPoolModel struct {
	tokens map[int64]*accessTokenModel
	sync.RWMutex
}

type wecomAPIResponse struct {
	ErrCode        int64  `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	Invaliduser    string `json:"invaliduser"`
	Invalidparty   string `json:"invalidparty"`
	Invalidtag     string `json:"invalidtag"`
	Unlicenseduser string `json:"unlicenseduser"`
	Msgid          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}

var tokenPool = tokenPoolModel{
	tokens:  make(map[int64]*accessTokenModel),
	RWMutex: sync.RWMutex{},
}

func getAccessToken(secret *ModWecomModel, enforce bool) *accessTokenModel {
	tokenPool.Lock()
	defer tokenPool.Unlock()
	token, ok := tokenPool.tokens[secret.ID]
	if ok && !enforce {
		if time.Now().Unix() < token.ExpiresIn {
			return token
		}
	}
	client := &http.Client{}
	req, err := client.Get(fmt.Sprintf(getTokenUrl, secret.CorpID, secret.Secret))
	if err != nil {
		glog.Warning("failed to get token [%s]", err.Error())
		return nil
	}
	defer func() {
		_ = req.Body.Close()
	}()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		glog.Warning("failed to read body [%s]", err.Error())
		return nil
	}

	ac := new(accessTokenModel)
	err = json.Unmarshal(body, ac)
	if err != nil {
		glog.Warning("failed to unmarshal [%s]", err.Error())
		return nil
	}
	ac.ExpiresIn = time.Now().Unix() + ac.ExpiresIn
	tokenPool.tokens[secret.ID] = ac
	return ac
}
