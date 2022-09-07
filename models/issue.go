package models

import (
	"altsub/base"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MIssueHandling struct {
	TX          *gorm.DB          `json:"-" gorm:"-"`
	ID          uint              `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time         `json:"created_at"`
	DeletedAt   gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
	Description string            `json:"description" gorm:"column:col_description;type:text"`
	Last        int64             `json:"last" gorm:"column:col_last;comment:预估处理时长"`
	Username    string            `json:"username" gorm:"column:col_username;not null;comment:告警认领人"`
	User        MUser             `json:"user" gorm:"foreignKey:Username;references:Username"`
	EventId     string            `json:"eventid" gorm:"column:col_eventid;not null;comment:认领告警id"`
	Events      []*MSchemaedEvent `json:"events" gorm:"-"`
}

type MIssueHandlings struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MIssueHandling
}

func (*MIssueHandling) TableName() string {
	return "tb_issues"
}

func (i *MIssueHandling) Add() error {
	if i.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增问题处理项失败", "models:issue.Add()")
		return err
	}
	return i.TX.Create(i).Error
}

func (i *MIssueHandling) SoftDelete() error {
	if i.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "软删除问题处理项失败", "models:issue.SoftDelete()")
		return err
	}
	return i.TX.Delete(i).Error
}

func (i *MIssueHandling) HardDelete() error {
	if i.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "硬删除问题处理项失败", "models:issue.HardDelete()")
		return err
	}
	return i.TX.Unscoped().Delete(i).Error
}

func (i *MIssueHandling) Update() error {
	if i.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "更新问题处理项失败", "models:issue.Update()")
		return err
	}
	return i.TX.Save(i).Error
}

func (i *MIssueHandling) Get(preloads ...string) error {
	if i.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取问题处理项失败", "models:issue.Get()")
		return err
	}
	for _, p := range preloads {
		i.TX = i.TX.Preload(p)
	}
	return i.TX.First(i).Error
}

func (is *MIssueHandlings) Fetch(preloads ...string) error {
	if is.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取问题处理项列表失败", "models:issue.Fetch()")
		return err
	}
	for _, p := range preloads {
		is.TX = is.TX.Preload(p)
	}
	return is.PQ.Query(is.TX, &is.All, &MIssueHandling{}).Error
}

func (i *MIssueHandling) FetchEvents() error {
	if i.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取此次告警事件列表失败", "models:issue.FetchEvents()")
		return err
	}
	if len(i.EventId) <= 0 {
		err := errors.New("empty eventid")
		base.NewLog("error", err, "获取此次告警事件列表失败", "models:issue.FetchEvents()")
		return err
	}
	var ses = MSchemaedEvents{
		TX: i.TX,
		PQ: PageQuery{
			Size:   100,
			Page:   1,
			Order:  "-created_at",
			Search: fmt.Sprintf("col_eventid: %s", i.EventId),
		},
	}
	if err := ses.FetchAfter(time.Now().Local().Add(-1 * 7 * 24 * time.Hour)); err != nil {
		base.NewLog("error", err, "获取此次告警事件列表失败", "models:issue.FetchEvents()")
		return err
	}
	i.Events = ses.All
	return nil
}
