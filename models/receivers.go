package models

import (
	"altsub/base"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	ReceiverTypeDingtalkApp = iota + 1
	ReceiverTypeDingtalkPersonal
	ReceiverTypeSMS
	ReceiverTypeVoice
	ReceiverTypeDingtalkRobot
)

type MReceiver struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index" `
	TX          *gorm.DB       `json:"-" gorm:"-"`
	Type        int            `json:"type" gorm:"column:col_type;not null;default:1;uniqueIndex:idx_name_type;comment:告警接收者类型标志"`
	Name        string         `json:"name" gorm:"column:col_name;type:varchar(64);not null;uniqueIndex:idx_name_type;comment:告警接收者名称"`
	Description string         `json:"description" gorm:"column:col_description;type:text;comment:告警接收者描述信息"`
	Auth        JSON           `json:"auth" gorm:"column:col_auth;not null;comment:告警接收者认证信息"`
	AuthHash    string         `json:"auth_hash" gorm:"column:col_auth_hash;type:varchar(64);not null;uniqueIndex;comment:认证信息hash"`
	Rules       []MRule        `json:"rules" gorm:"many2many:tb_subscribes;foreignKey:ID;joinForeignKey:col_receiver_id;joinReferences:col_rule_id;constraint:OnDelete:CASCADE;"`
}

type MReceivers struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MReceiver
}

func (*MReceiver) TableName() string {
	return "tb_receivers"
}

// func (r *MReceiver) AfterFind(tx *gorm.DB) error {
// 	r.Auth = nil
// 	return nil
// }

func (r *MReceiver) BeforeCreate(tx *gorm.DB) error {
	_m := map[string]interface{}{}
	_b := []byte{}
	json.Unmarshal(r.Auth, &_m)
	_b, _ = json.Marshal(_m)
	hs := md5.New()
	hs.Write(_b)
	r.AuthHash = fmt.Sprintf("%x", hs.Sum(nil))
	return nil
}

func (rcvs *MReceivers) Fetch(preloads ...string) error {
	if rcvs.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取接收者失败", "models:receiver.Fetch()")
		return err
	}
	for _, p := range preloads {
		rcvs.TX = rcvs.TX.Preload(p)
	}
	if err := rcvs.PQ.Query(rcvs.TX, &rcvs.All, &MReceiver{}).Error; err != nil {
		base.NewLog("error", err, "获取接收者", "models:receiver.Fetch()")
		return err
	}
	return nil
}

func (r *MReceiver) Get(preloads ...string) error {
	for _, p := range preloads {
		r.TX = r.TX.Preload(p)
	}
	return r.TX.First(r).Error
}

func (r *MReceiver) Delete() error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "删除接收者失败", "models:receiver.Delete()")
		return err
	}
	if err := r.TX.Unscoped().Delete(r).Error; err != nil {
		base.NewLog("error", err, "删除接收者失败", "models:receiver.Delete()")
		return err
	}
	return nil
}

func (r *MReceiver) Add() error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增接收者失败", "models:receiver.Add()")
		return err
	}
	if len(r.Auth) <= 0 {
		err := errors.New("empty auth info")
		base.NewLog("error", err, "新增接收者失败", "models:receiver.Add()")
		return err
	}
	if r.Type <= 0 || r.Type > ReceiverTypeDingtalkRobot {
		err := errors.New("wrong receiver type")
		base.NewLog("error", err, "新增接收者失败", "models:receiver.Add()")
		return err
	}
	if err := r.TX.Create(r).Error; err != nil {
		base.NewLog("error", err, "新增接收者", "models:receiver.Add()")
		return err
	}
	return nil
}

func (rcv *MReceiver) Subscribe(rs []MRule) error {
	if rcv.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "订阅失败", "models:receiver.Subscribe()")
		return err
	}
	rcv.Rules = rs
	if err := rcv.TX.Select("Rules").Save(rcv).Error; err != nil {
		base.NewLog("error", err, "订阅规则", "models:receiver.Subscribe()")
		return err
	}
	return nil
}
