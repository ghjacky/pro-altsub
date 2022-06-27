package receiver

import (
	"altsub/base"
	"altsub/models"
	"encoding/json"
	"errors"
)

func Add(rcv *models.MReceiver) error {
	if rcv.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增接收者失败", "receiver:Add()")
		return err
	}
	_v := map[string]interface{}{}
	if err := json.Unmarshal(rcv.Auth, &_v); err != nil {
		base.NewLog("error", err, "新增接收者失败", "receiver:Add()")
		return err
	}
	return rcv.Add()
}

func Fetch(rcvs *models.MReceivers) error {
	if rcvs.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取接收者失败", "receiver:Fetch()")
		return err
	}
	return rcvs.Fetch()
}

func Get(rcv *models.MReceiver) error {
	if rcv.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取接收者失败", "receiver:Get()")
		return err
	}
	return rcv.Get("Rules")
}
