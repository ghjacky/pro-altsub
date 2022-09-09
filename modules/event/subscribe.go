package event

import (
	"altsub/base"
	"altsub/models"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkSubscribe(rs []models.MRule) []models.MReceiver {
	var rcvs = []models.MReceiver{}
	for _, r := range rs {
		rcvs = append(rcvs, r.Receivers...)
	}
	return rcvs
}

func getDefaultServiceGroup(bigtype, service string) []models.MReceiver {
	var rcvs = []models.MReceiver{}
	hc := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	hq, _ := http.NewRequest("GET", fmt.Sprintf("%s?big_type=%s&env=%s&service=%s", base.Config.MainConfig.Easy, bigtype, "online", service), nil)
	res, err := hc.Do(hq)
	if err != nil {
		base.NewLog("error", err, "easy接口请求失败", "getDefaultServiceGroup()")
		return rcvs
	}
	defer res.Body.Close()
	var _res = map[string]interface{}{}
	d, e := ioutil.ReadAll(res.Body)
	if e != nil {
		base.NewLog("error", err, "easy接口响应体读取失败", "getDefaultServiceGroup()")
		return rcvs
	}
	e = json.Unmarshal(d, &_res)
	if e != nil {
		base.NewLog("error", err, "easy接口响应体解析失败", "getDefaultServiceGroup()")
		return rcvs
	}
	data, _ := _res["data"].(map[string]interface{})
	if data == nil {
		base.NewLog("error", err, "easy接口响应体字段缺失", "getDefaultServiceGroup()")
		return rcvs
	}
	chatid, _ := data["group_dd_token"].(string)
	if len(chatid) == 0 {
		return rcvs
	}
	auth := map[string]interface{}{
		"app_key":    base.Config.NotificationConf.DefaultDingTalkAppKey,
		"app_secret": base.Config.NotificationConf.DefaultDingTalkAppSecret,
		"chat_id":    chatid,
	}
	ab, _ := json.Marshal(auth)
	hs := md5.New()
	hs.Write(ab)
	rcvs = append(rcvs, models.MReceiver{
		Type:     models.ReceiverTypeDingtalkApp,
		Auth:     ab,
		AuthHash: fmt.Sprintf("%x", hs.Sum(nil)),
	})
	return rcvs
}
