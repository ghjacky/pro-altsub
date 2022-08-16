package models

import (
	"altsub/base"
	"errors"
	"time"

	"gorm.io/gorm"
)

////// 用户
type MUser struct {
	Username string       `json:"username" gorm:"type:varchar(64);column:col_username;primaryKey"`
	Phone    string       `json:"phone" gorm:"type:varchar(16);column:col_phone"`
	Nickname string       `json:"nickname" gorm:"column:col_nickname"`
	Fullname string       `json:"fullname" gorm:"column:col_fullname"`
	Email    string       `json:"email" gorm:"column:col_email"`
	Position string       `json:"position" gorm:"column:col_position"`
	Groups   []MDutyGroup `json:"groups" gorm:"many2many:tb_duty_groups_users"`
}

type MUsers struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MUser
}

func (MUser) TableName() string {
	return "tb_users"
}

/////// 排班时间
type MDutyAt struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" `
	TX        *gorm.DB       `json:"-" gorm:"-"`
	StartAt   int64          `json:"startAt" gorm:"column:col_start_at"`
	EndAt     int64          `json:"endAt" gorm:"column:col_end_at"`
	GroupID   uint           `json:"groupId" gorm:"column:col_group_id"`
	// Type      uint8          `json:"type" gorm:"column:col_type;comment:"`
	// Year      uint16         `json:"year" gorm:"column:col_year;not null"`
	// Month     uint8          `json:"month" gorm:"column:col_month;not null"`
	// Day       JSON           `json:"day" gorm:"type:blob;column:col_day;not null"`
}

type MDutyAts struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MDutyAt
}

func (MDutyAt) TableName() string {
	return "tb_duty_ats"
}

////// 排班分组
type MDutyGroup struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	TX        *gorm.DB       `json:"-" gorm:"-"`
	Name      string         `json:"name" gorm:"column:col_name"`
	DutyID    uint           `json:"dutyId" gorm:"column:col_duty_id;not null"`
	Ats       []MDutyAt      `json:"ats" gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Users     []MUser        `json:"users" gorm:"many2many:tb_duty_groups_users"`
}

type MDutyGroups struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MDutyGroup
}

func (MDutyGroup) TableName() string {
	return "tb_duty_groups"
}

////// 排班班次
type MDuty struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" `
	TX        *gorm.DB       `json:"-" gorm:"-"`
	Name      string         `json:"name" gorm:"type:varchar(64);column:col_name;not null;uniqueIndex"`
	Groups    []MDutyGroup   `json:"groups" gorm:"foreignKey:DutyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// RuleID    uint           `json:"ruleId" gorm:"column:col_rule_id"`
	// Rule      MRule          `json:"rules" gorm:"foreignKey:RuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type MDuties struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MDuty
}

func (MDuty) TableName() string {
	return "tb_duties"
}

func (d *MDuty) Add() error {
	if d.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "增加排班班次失败", "models:Duty.Add()")
		return err
	}
	return d.TX.Save(d).Error
}

func (d *MDuty) HardDelete() error {
	if d.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "删除排班班次失败", "models:Duty.HardDelete()")
		return err
	}
	return d.TX.Unscoped().Delete(d).Error
}

func (d *MDuty) SoftDelete() error {
	if d.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "删除排班班次失败", "models:Duty.SoftDelete()")
		return err
	}
	return d.TX.Delete(d).Error
}

func (ds *MDuties) Fetch(preloads ...string) error {
	if ds.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "获取排班列表失败", "models:Duties.Fetch()")
		return err
	}
	for _, p := range preloads {
		ds.TX = ds.TX.Preload(p)
	}
	return ds.PQ.Query(ds.TX, &ds.All).Error
}

func (d *MDuty) AddGroup(g *MDutyGroup) error {
	if d.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "新增排班分组失败", "models:Duty.AddGroup()")
		return err
	}
	return d.TX.Model(d).Association("Groups").Append(g)
}

func (d *MDuty) DeleteGroup(g *MDutyGroup) error {
	if d.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "删除排班分组失败", "models:Duty.DeleteGroup()")
		return err
	}
	return d.TX.Model(d).Association("Groups").Delete(g)
}

/////// 问题处理
const (
	IssueStatusFiring uint8 = iota + 1
)

type MIssueHandling struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index" `
	TX          *gorm.DB       `json:"-" gorm:"-"`
	Title       string         `json:"title" gorm:"column:col_title;not null"`
	Description string         `json:"description" gorm:"type:text;column:col_description"`
	Status      uint8          `json:"status" gorm:"column:col_status;not null;default:1;comment:问题处理状态"`
	Username    string         `json:"username" gorm:"column:col_username"`
	User        MUser          `json:"user" gorm:"foreignKey:Username"`
}

type MIssueHandlings struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MIssueHandling
}

func (MIssueHandling) TableName() string {
	return "tb_issue_handlings"
}

func (ih *MIssueHandling) Add() error {
	if ih.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "新增问题处理失败", "models:MIssueHandling.Add()")
		return err
	}
	return ih.TX.Create(ih).Error
}

func (ih *MIssueHandling) HardDelete() error {
	if ih.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "删除问题处理失败", "models:MIssueHandling.Delete()")
		return err
	}
	return ih.TX.Unscoped().Delete(ih).Error
}

func (ih *MIssueHandling) Update() error {
	if ih.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "变更问题处理失败", "models:MIssueHandling.Update()")
		return err
	}
	return ih.TX.Save(ih).Error
}

func (ih *MIssueHandling) Get() error {
	if ih.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "获取问题处理失败", "models:MIssueHandling.Get()")
		return err
	}
	return ih.TX.First(ih).Error
}

func (ihs *MIssueHandlings) Fetch(preloads ...string) error {
	if ihs.TX == nil {
		err := errors.New("nil db error")
		base.NewLog("error", err, "获取问题处理列表失败", "models:MIssueHandling.Fetch()")
		return err
	}
	for _, p := range preloads {
		ihs.TX = ihs.TX.Preload(p)
	}
	return ihs.PQ.Query(ihs.TX, &ihs.All).Error
}
