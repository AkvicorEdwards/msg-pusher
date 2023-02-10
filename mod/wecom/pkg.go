package wecom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AkvicorEdwards/glog"
	"io/ioutil"
	"net/http"
)

type Package struct {
	wecom *ModWecomModel
	msg   []byte
}

func (p *Package) Send() bool {
	actoken := getAccessToken(p.wecom, false)
	if actoken == nil {
		glog.Debug("failed to get token")
		return false
	}
	url := fmt.Sprintf(sendUrl, actoken.AccessToken)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(p.msg))
	if err != nil {
		glog.Debug("failed to make request [%s]", err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	for err != nil {
		glog.Debug("failed to send request [%s]", err.Error())
		return false
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Debug("failed to read body [%s]", err.Error())
		return false
	}
	var res = &wecomAPIResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		glog.Debug("failed to unmarshal body [%s]", err.Error())
		return false
	}
	if res.ErrCode != 0 {
		glog.Warning("failed to send [%#v]", res)
		return false
	}
	return true
}
