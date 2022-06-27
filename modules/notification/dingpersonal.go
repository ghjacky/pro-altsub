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
	DingPersonalMsgTpl = DingTalkMsgTpl
)

type DingPersonal struct {
	Mobiles string
	Client  *http.Client
	Event   schema.SchemaedEvent
}

func (dp *DingPersonal) NewRequest(mobiles, content string) *http.Request {
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
	db, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://tnotice.songguo7.com/tnotice/dingtalk/send", bytes.NewReader(db))
	if err != nil {
		base.NewLog("error", err, "dingtalk工作通知请求创建失败", "notification:dingpersonal.NewRequest()")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (dp *DingPersonal) Notice(msg string) error {
	base.NewLog("info", nil, fmt.Sprintf("dingtalk工作通知发送：%s", dp.Mobiles), "notification:dingpersonal.Notice()")
	resp, err := dp.Client.Do(dp.NewRequest(dp.Mobiles, msg))
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()
		var result = map[string]interface{}{}
		b, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(b, &result); err != nil {
			base.NewLog("error", err, "dingtalk工作通知发送失败", "notification:dingpersonal.Notice()")
			return err
		} else {
			base.NewLog("debug", nil, fmt.Sprintf("dingtalk工作通知发送失败：%s", string(b)), "notification:dingpersonal.Notice()")
		}
	}
	return nil
}

func (dp *DingPersonal) SetEvent(ev schema.SchemaedEvent) {
	dp.Event = ev
}

func (dt *DingPersonal) RenderMsg() string {
	tpl, err := template.New("alert").Parse(DingPersonalMsgTpl)
	if err != nil {
		base.NewLog("error", err, "dingtalk工作通知模板创建失败", "notification:dingpersonal.RenderMsg()")
		return ""
	}
	var msg = bytes.Buffer{}
	if err := tpl.Execute(&msg, dt.Event); err != nil {
		base.NewLog("error", err, "dingtalk工作通知模板渲染失败", "notification:dingpersonal.RenderMsg()")
		return ""
	}
	return msg.String()
}

func (dp *DingPersonal) ParseAuth(auth []byte) {
	var au = map[string]string{}
	if err := json.Unmarshal(auth, &au); err != nil {
		base.NewLog("error", err, "dingpersonal认证信息解析失败", "notification:dingpersonal.ParseAuth()")
		return
	}
	dp.Mobiles = string(au["mobiles"])
}
