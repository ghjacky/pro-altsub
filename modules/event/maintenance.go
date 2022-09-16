package event

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/maintenance"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkMaintenance(rs []models.MRule) bool {
	return maintenance.Check(rs)
}

func checkPublish(bigtype, service, instance string) bool {
	hc := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	base.NewLog("debug", nil, fmt.Sprintf("请求easy接口检测服务发版状态信息（%s/api_deploy_check?big_type=%s&env=%s&service=%s&instance=%s)", base.Config.MainConfig.Easy, bigtype, "online", service, instance), "checkPublish()")
	hq, _ := http.NewRequest("GET", fmt.Sprintf("%s/api_deploy_check?big_type=%s&env=%s&service=%s&instance=%s", base.Config.MainConfig.Easy, bigtype, "online", service, instance), nil)
	res, err := hc.Do(hq)
	if err != nil {
		base.NewLog("error", err, "easy接口请求失败", "checkMaintenance()")
		return false
	}
	defer res.Body.Close()
	var _res = map[string]interface{}{}
	d, e := ioutil.ReadAll(res.Body)
	if e != nil {
		base.NewLog("error", err, "easy接口响应体读取失败", "checkMaintenance()")
		return false
	}
	e = json.Unmarshal(d, &_res)
	if e != nil {
		base.NewLog("error", err, "easy接口响应体解析失败", "checkMaintenance()")
		return false
	}
	base.NewLog("debug", nil, fmt.Sprintf("从easy检测服务发版状态返回结果：%v", _res), "checkPublish()")
	v, ok := _res["data"].(bool)
	if !v || !ok {
		return false
	}
	return true
}
