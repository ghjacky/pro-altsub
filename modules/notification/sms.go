package notification

import (
	"altsub/base"
	"altsub/modules/schema"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	SMSMsgTpl = DingTalkMsgTpl
)

type SMS struct {
	Mobiles string
	Client  *http.Client
	Event   schema.SchemaedEvent
}

func (s *SMS) NewRequest(mobiles, content string) *http.Request {
	var data = map[string]interface{}{}
	var secret = "ab37655347122540857b44950fca0ba3"
	data["timestamp"] = fmt.Sprintf("%d", time.Now().Local().Unix())
	data["rand"] = fmt.Sprintf("%d", rand.Intn(100000-1000)+1000)
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s%s%s", secret, data["timestamp"], data["rand"])))
	data["sign"] = fmt.Sprintf("%x", h.Sum(nil))
	data["appid"] = "520680"
	data["from"] = "ebike-order-core"
	data["mobiles"] = mobiles
	data["content"] = content
	data["channel"] = 1
	base.NewLog("debug", nil, "创建短信接口请求", "notification:sms.NewRequest()")
	db, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://tnotice.songguo7.com/tnotice/sms/sendAlarm", bytes.NewReader(db))
	if err != nil {
		base.NewLog("error", err, "创建短信接口请求失败", "notification:sms.NewRequest()")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (s *SMS) Notice(msg string) error {
	base.NewLog("info", nil, fmt.Sprintf("告警短信发送：%s", s.Mobiles), "notification:sms.Notice()")
	resp, err := s.Client.Do(s.NewRequest(s.Mobiles, msg))
	if err != nil {
		base.NewLog("error", err, "告警短信发送失败", "notification:sms.Notice()")
		return err
	} else {
		var result = map[string]interface{}{}
		b, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(b, &result); err != nil {
			base.NewLog("error", err, fmt.Sprintf("告警短信发送失败：%s", string(b)), "notification:sms.Notice()")
			return err
		} else {
			base.NewLog("debug", nil, fmt.Sprintf("告警短信发送响应：%s", string(b)), "notification:sms.Notice()")
		}
	}
	return nil
}

func (s *SMS) RenderMsg() string {
	tpl, err := template.New("sms_tpl").Parse(SMSMsgTpl)
	if err != nil {
		base.NewLog("error", err, "短信告警模板创建失败", "notification:sms.RenderMsg()")
		return ""
	}
	var msg = bytes.Buffer{}
	if err := tpl.Execute(&msg, s.Event); err != nil {
		base.NewLog("error", err, "短信告警模板渲染失败", "notification:sms.RenderMsg()")
		return ""
	}
	return msg.String()
}

func (s *SMS) ParseAuth(auth []byte) {
	var au = map[string]string{}
	if err := json.Unmarshal(auth, &au); err != nil {
		base.NewLog("error", err, "短信接口认证信息解析失败", "notification:sms.ParseAuth()")
		return
	}
	s.Mobiles = string(au["mobiles"])
}

func (s *SMS) SetEvent(ev schema.SchemaedEvent) {
	s.Event = ev
}
