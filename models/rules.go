package models

import (
	"altsub/base"
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	RuleTypeDefault = iota
	RuleTypeSubscribe
	RuleTypeEventRelationship // 告警认领（同一类告警）
	RuleTypeMaintenance
	RuleTypeSuppression
	RuleLogicIntersection = iota - 1
	RuleLogicConcatenation
)

// 规则：（关联：规则、接收者、维护项）
type MRule struct {
	// StartAt     int64          `json:"startAt" gorm:"column:col_start_at;not null;comment:规则生效时间段开始秒级时间戳"`
	// EndAt       int64          `json:"end_at" gorm:"column:col_end_at;not null;comment:规则生效时间段结束秒级时间戳"`
	//（belongs_to、has_one关系中，如果新增条目时要使外键为空，则须使用*uint类型，因为*类型零值为nil，而uint零值则为0，如果是0而库认为对应id为0的关联项不存在，数据库会报错）
	Type      uint      `json:"type" gorm:"column:col_type;not null;comment:规则类型（订阅、维护、抑制）"`
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	// DeletedAt    gorm.DeletedAt   `json:"deleted_at" gorm:"index" `
	TX           *gorm.DB         `json:"-" gorm:"-"`
	Name         string           `json:"name" gorm:"column:col_name;type:varchar(64);not null;uniqueIndex:idx_source_rule;comment:规则名称"`
	Description  string           `json:"description" gorm:"column:col_description;type:text;comment:规则描述信息"`
	PrevID       *uint            `json:"prev_id" gorm:"column:col_prev_id;constraint:OnUpdate:CASCADE,ONDELETE:SET NULL;comment:规则链条下一节点id"`
	Prev         *MRule           `json:"prev" gorm:"foreignKey:PrevID;references:ID;constraint:OnUpdate:CASCADE,ONDELETE:SET NULL"`
	Logic        int              `json:"logic" gorm:"column:col_logic;not null;default:1;comment:规则内容中每项之间的逻辑关系（且、或）"`
	SourceID     uint             `json:"source_id" gorm:"column:col_source_id;not null;uniqueIndex:idx_source_rule;comment:schema相关联的 source id"`
	Source       *MSource         `json:"source" gorm:"foreignKey:SourceID;references:ID"`
	Maintenances []MMaintenance   `json:"maintenances" gorm:"many2many:tb_maintenances_rules;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Clauses      []MRuleClause    `json:"clauses" gorm:"foreignKey:RuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Receivers    []MReceiver      `json:"receivers" gorm:"many2many:tb_subscribes;references:ID;joinForeignKey:col_rule_id;joinReferences:col_receiver_id;constraint:OnDelete:CASCADE;"`
	Events       []MSchemaedEvent `json:"events" gorm:"many2many:tb_events_rules"`
}

type MRules struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []MRule
}

func (*MRule) TableName() string {
	return "tb_rules"
}

func (r *MRule) Get(preloads ...string) error {
	for _, p := range preloads {
		r.TX = r.TX.Preload(p)
	}
	return r.TX.First(r).Error
}

func (rs *MRules) Fetch(preloads ...string) error {
	if rs.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则失败", "models:rule.Fetch()")
		return err
	}
	for _, p := range preloads {
		rs.TX = rs.TX.Preload(p)
	}
	if err := rs.PQ.Query(rs.TX, &rs.All, &MRule{}).Error; err != nil {
		base.NewLog("error", err, "获取规则", "models:rule.Fetch()")
		return err
	}
	return nil
}

func (r *MRule) Add() error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增规则失败", "models:rule.Add()")
		return err
	}

	if len(r.Name) == 0 {
		err := errors.New("empty rule name")
		base.NewLog("error", err, "新增规则失败", "models:rule.Add()")
		return err
	}

	if r.Logic != RuleLogicConcatenation && r.Logic != RuleLogicIntersection {
		err := errors.New("wrong rule logic operator")
		base.NewLog("error", err, "新增规则失败", "models:rule.Add()")
		return err
	}

	if r.Source.ID == 0 {
		if len(r.Source.Name) == 0 {
			err := errors.New("wrong ralated source")
			base.NewLog("error", err, "新增规则失败", "models:rule.Add()")
			return err
		} else {
			if err := r.Source.GetByName(); err != nil {
				base.NewLog("error", err, "新增规则失败", "models:rule.Add()")
				return err
			}
		}
	}

	if err := r.TX.Create(r).Error; err != nil {
		base.NewLog("error", err, "新增规则", "models:rule.Add()")
		return err
	}
	return nil
}

func (r *MRule) GetByNameAndSource() error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则失败", "models:rule.GetByNameAndSource()")
		return err
	}
	if len(r.Name) <= 0 || r.SourceID == 0 {
		err := errors.New("empty rule name or source id")
		base.NewLog("error", err, "获取规则失败", "models:rule.GetByNameAndSource()")
		return err
	}
	if err := r.TX.Where("col_name = ? and col_source_id = ?", r.Name, r.SourceID).First(r).Error; err != nil {
		base.NewLog("error", err, "获取规则", "models:rule.GetByNameAndSource()")
		return err
	}
	return nil
}

func (r *MRule) GetChain(preloads ...string) error {
	db := base.DB()
	if db == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则链", "models:rule.GetChain()")
		return err
	}
	if r.ID == 0 {
		if len(r.Name) <= 0 || r.SourceID == 0 {
			err := errors.New("empty id and (name or source id) of rule")
			base.NewLog("error", err, "获取规则链失败", "models:rule.GetChain()")
			return err
		} else {
			if err := r.GetByNameAndSource(); err != nil {
				base.NewLog("error", err, "获取规则链失败", "models:rule.GetChain()")
				return err
			}
		}
	}
	for _, prl := range preloads {
		db = db.Preload(prl)
	}
	if err := db.First(r).Error; err != nil {
		base.NewLog("error", err, "获取规则链失败", "models:rule.GetChain()")
		return err
	} else {
		if r.PrevID == nil {
			return nil
		} else {
			if r.Prev == nil {
				err := errors.New("no preloading prev")
				base.NewLog("error", err, "获取规则链失败", "models:rule.GetChain()")
				return err
			}
			r.Prev.TX = db
			if err := r.Prev.GetChain("Prev", "Clauses"); err != nil {
				base.NewLog("error", err, "获取规则链", "models:rule.GetChain()")
				return err
			}
			return nil
		}
	}
}

func (r *MRule) Assign(rcvs []MReceiver) error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "指派失败", "models:rule.Assign()")
		return err
	}

	if err := r.TX.Model(r).Association("Receivers").Append(&rcvs); err != nil {
		base.NewLog("error", err, "指派规则失败", "models:rule.Assign()")
		return err
	}
	return nil
}

func (r *MRule) Delete() error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "规则删除失败", "models:rule.Delete()")
		return err
	}
	var s = &MSubscribe{RuleID: r.ID}
	if err := r.TX.Unscoped().Delete(r).Error; err != nil {
		base.NewLog("error", err, "规则删除失败", "models:rule.Delete()")
		return err
	}
	s.DeleteByRule(r.TX)
	return nil
}
