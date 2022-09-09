package notification

import (
	"altsub/base"
	"altsub/modules/schema"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

const (
	DingTalkChatApi = "https://oapi.dingtalk.com/chat/send"
	DingTalkMsgTpl  = `
	{{- range $v := . -}}
		{{- if eq $v.Key "status" }}
			{{- if eq $v.Value "firing" }}
# **🆘 故障告警**
--------------------------------------------------
			{{- else if eq $v.Value "resolved" }}
# **✅ 故障恢复**
--------------------------------------------------
			{{- end}}
		{{- else }}
			{{- with $v.CName }}
				{{- if and (ne . "") ($v.Value) }}
					{{- if eq $v.SType "text" }}
- **{{ . }}:** {{ $v.Value }}
					{{- else if eq $v.SType "link" }}
- **点击查看: [{{ . }}]({{ $v.Value }})**
					{{- else }}
- **{{ . }}:** {{ $v.Value }}
					{{- end }}
				{{- end }}
			{{- end }}
		{{- end }}
	{{- end }}
--------------------------------------------------
	`
)

type DingTalk struct {
	Token  string
	ChatID string
	Client *http.Client
	Event  schema.SchemaedEvent
}

func (dt *DingTalk) SetEvent(ev schema.SchemaedEvent) {
	dt.Event = ev
}

func (dt *DingTalk) RenderMsg() string {
	tpl, err := template.New("alert").Parse(DingTalkMsgTpl)
	if err != nil {
		base.NewLog("error", err, "dingtalk告警模板创建失败", "notification:dingtalk.RenderMsg()")
		return ""
	}
	var msg = bytes.Buffer{}
	if err := tpl.Execute(&msg, dt.Event); err != nil {
		base.NewLog("error", err, "dingtalk告警模板渲染失败", "notification:dingtalk.RenderMsg()")
		return ""
	}
	return msg.String()
}

func (dt *DingTalk) Notice(msg string) error {
	base.NewLog("info", nil, fmt.Sprintf("dingtalk告警发送：%s", dt.ChatID), "notification:dingtalk.Notice()")
	url := strings.Builder{}
	url.WriteString(DingTalkChatApi)
	url.WriteString("?access_token=")
	url.WriteString(dt.Token)
	bm := dt.makeActionCardMsgContent(msg)
	bb, _ := json.Marshal(bm)
	body := strings.NewReader(string(bb))
	req, err := http.NewRequest("POST", url.String(), body)
	if err != nil {
		base.NewLog("error", err, "dingtalk接口请求创建失败", "notification:dingtalk.Notice()")
		return err
	}
	resp, err := dt.Client.Do(req)
	if err != nil {
		base.NewLog("error", err, "dingtalk接口请求失败", "notification:dingtalk.Notice()")
		return err
	}
	defer resp.Body.Close()
	var result = map[string]interface{}{}
	rb, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(rb, &result); err != nil {
		base.NewLog("error", err, "dingtalk告警发送失败", "notification:dingtalk.Notice()")
		return err
	} else if result["errcode"] != 0 {
		errmsg, _ := result["errmsg"].(string)
		if errmsg == "ok" {
			return nil
		} else {
			err := errors.New(errmsg)
			base.NewLog("error", err, "dingtalk告警发送失败", "notification:dingtalk.Notice()")
			return err
		}
	}
	return nil
}

func (dt *DingTalk) ParseAuth(auth []byte) {
	var au = struct {
		AppKey    string `json:"app_key"`
		AppSecret string `json:"app_secret"`
		ChatId    string `json:"chat_id"`
	}{}
	if err := json.Unmarshal(auth, &au); err != nil {
		base.NewLog("error", err, "dingtalk认证信息解析失败", "notification:dingtalk.ParseAuth()")
	}
	if len(au.AppKey) == 0 || len(au.AppSecret) == 0 {
		au.AppKey = base.Config.DefaultDingTalkAppKey
		au.AppSecret = base.Config.DefaultDingTalkAppSecret
	}
	if len(au.ChatId) == 0 {
		au.ChatId = base.Config.DefaultDingtalkChatID
	}
	dt.Token = GetDingTalkChatAccessToken(au.AppKey, au.AppSecret)
	dt.ChatID = au.ChatId
}

func GetDingTalkChatAccessToken(appKey, appSecret string) string {
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", appKey, appSecret)
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		base.NewLog("error", err, "dingtalk获取token接口请求创建失败", "notification:GetDingTalkChatAccessToken()")
		return ""
	}
	resp, err := hc.Do(req)
	if err != nil {
		base.NewLog("error", err, "dingtalk获取token接口请求失败", "notification:GetDingTalkChatAccessToken()")
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		rb, _ := ioutil.ReadAll(resp.Body)
		rm := map[string]interface{}{}
		json.Unmarshal(rb, &rm)
		return rm["access_token"].(string)
	}
	base.NewLog("error", nil, "dingtalk获取token接口请求失败", "notification:GetDingTalkChatAccessToken()")
	return ""
}

func (dt *DingTalk) makeActionCardMsgContent(msg string) map[string]interface{} {

	// ensureUrl := fmt.Sprintf("https://oapi.dingtalk.com/connect/oauth2/sns_authorize?appid=%s&response_type=code&scope=snsapi_auth&state=%s&redirect_uri=%s", base.Config.NotificationConf.DefaultDingTalkAppKey, generateEventEnsureQueryString(dt.Event), fmt.Sprintf("%s%s", base.Config.MainConfig.Domain, base.Config.NotificationConf.AlertEnsureApi))

	// dissUrl := ""
	var status = "firing"
	for _, ei := range dt.Event {
		if ei.Key == "status" {
			status, _ = ei.Value.(string)
		}
	}
	if status == "firing" {
		return map[string]interface{}{
			"chatid": dt.ChatID,
			"msg": map[string]interface{}{
				"msgtype": "action_card",
				"action_card": map[string]interface{}{
					"title":           "松果告警",
					"markdown":        msg,
					"btn_orientation": 1, // 1: 按钮横向排列； 0: 竖向排列
					"btn_json_list": []map[string]string{
						{
							"title":      "告警确认",
							"action_url": "ensureUrl",
						},
						// {
						// 	"title":      "无效告警",
						// 	"action_url": dissUrl,
						// },
					},
				},
			},
		}
	} else {
		return map[string]interface{}{
			"chatid": dt.ChatID,
			"msg": map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]interface{}{
					"title": "松果告警",
					"text":  msg,
				},
			},
		}
	}

}

// func generateEventEnsureQueryString(ev schema.SchemaedEvent) string {
// 	return url.QueryEscape(fmt.Sprintf("event_id=%d&eventid=%s", ev.ID, ev.EventID))
// }
