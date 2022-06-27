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
	VoiceMsgTpl = DingTalkMsgTpl
)

type Voice struct {
	Mobiles string
	Client  *http.Client
	Event   schema.SchemaedEvent
}

func (v *Voice) NewRequest(mobiles, content string) *http.Request {
	var data = map[string]interface{}{}
	var secret = "ab37655347122540857b44950fca0ba3"
	data["timestamp"] = fmt.Sprintf("%d", time.Now().Local().Unix())
	data["rand"] = fmt.Sprintf("%d", rand.Intn(100000-1000)+1000)
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s%s%s", secret, data["timestamp"], data["rand"])))
	data["sign"] = fmt.Sprintf("%x", h.Sum(nil))
	data["appid"] = "520680"
	data["from"] = "ebike-order-core"
	data["alarmReceivers"] = mobiles
	data["channel"] = 1
	db, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://tnotice.songguo7.com/tnotice/call/callAlarm", bytes.NewReader(db))
	if err != nil {
		base.NewLog("error", err, "语音接口请求创建失败", "notification:voice.NewRequest()")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (v *Voice) Notice(msg string) error {
	base.NewLog("info", nil, fmt.Sprintf("语音告警发送：%s", v.Mobiles), "notification:voice.Notice")
	resp, err := v.Client.Do(v.NewRequest(v.Mobiles, msg))
	if err != nil {
		base.NewLog("error", err, "语音告警发送失败", "notification:voice.Notice()")
		return err
	} else {
		var result = map[string]interface{}{}
		b, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(b, &result); err != nil {
			base.NewLog("error", err, "语音告警发送失败", "notification:voice.Notice()")
			return err
		} else {
			base.NewLog("debug", nil, fmt.Sprintf("语音告警发送响应：%s", string(b)), "notification:voice.Notice()")
		}
	}
	return nil
}

func (v *Voice) RenderMsg() string {
	tpl, err := template.New("voice_tpl").Parse(SMSMsgTpl)
	if err != nil {
		base.NewLog("error", err, "语音告警模板创建失败", "notification:voice.RenderMsg()")
		return ""
	}
	var msg = bytes.Buffer{}
	if err := tpl.Execute(&msg, v.Event); err != nil {
		base.NewLog("error", err, "语音告警模板渲染失败", "notification:voice.RenderMsg()")
		return ""
	}
	return msg.String()
}

func (v *Voice) ParseAuth(auth []byte) {
	var au = map[string]string{}
	if err := json.Unmarshal(auth, &au); err != nil {
		base.NewLog("error", err, "语音认证信息解析失败", "notification:voice.ParseAuth()")
		return
	}
	v.Mobiles = string(au["mobiles"])
}

func (v *Voice) SetEvent(ev schema.SchemaedEvent) {
	v.Event = ev
}
