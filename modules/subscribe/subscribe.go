package subscribe

import (
	"altsub/base"
	"altsub/models"
	"errors"
)

func Subscribe(rcv models.MReceiver, rs []models.MRule) error {
	if rcv.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "订阅规则失败", "subscribe:Subescribe()")
		return err
	}
	if len(rs) <= 0 {
		err := errors.New("empty rules")
		base.NewLog("error", err, "订阅规则失败", "subscribe:Subescribe()")
		return err
	}
	return rcv.Subscribe(rs)
}

func Assign(r models.MRule, rcvs []models.MReceiver) error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "指派规则失败", "subscribe:Assign()")
		return err
	}
	if len(rcvs) <= 0 {
		err := errors.New("empty receivers")
		base.NewLog("error", err, "指派规则失败", "subscribe:Assign()")
		return err
	}
	return r.Assign(rcvs)
}

func Fetch(ss *models.MSubscribes) error {
	if ss.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取订阅/指派关系失败", "subscribe:Fetch()")
		return err
	}
	return ss.Fetch()
}
