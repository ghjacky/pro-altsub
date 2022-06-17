package models

import (
	"altsub/base"
	"errors"

	"gorm.io/gorm"
)

const (
	RuleTypeSubscribe = iota + 1
	RuleTypeEventRelationship
	RuleLogicIntersection = iota - 1
	RuleLogicConcatenation
	RuleOpEqual = iota - 3
	RuleOpGreatThan
	RuleOpLessThan
	RuleOpInclude
	RuleOpRegex
)

// 规则：（关联：规则、接收者、维护项）
type MRule struct {
	BaseModel
	Name     string   `json:"name" gorm:"column:col_name;type:varchar(64);not null;uniqueIndex;comment:规则名称"`
	StartAt  int64    `json:"start_at" gorm:"column:col_start_at;not null;comment:规则生效时间段开始秒级时间戳"`
	EndAt    int64    `json:"end_at" gorm:"column:col_end_at;not null;comment:规则生效时间段结束秒级时间戳"`
	PrevID   uint     `json:"prev_id" gorm:"column:col_prev_id;comment:规则链条下一节点id"`
	Prev     *MRule   `json:"prev" gorm:"foreignKey:PrevID;references:ID"`
	Logic    int      `json:"logic" gorm:"column:col_logic;not null;default:1;comment:规则内容中每项之间的逻辑关系（且、或）"`
	Operator int      `json:"operator" gorm:"column:col_operator;not null;default:1;comment:规则内容中每一项大比对符号"`
	SourceID uint     `json:"source_id" gorm:"column:col_source_id;not null;comment:schema相关联的 source id"`
	Source   *MSource `json:"source" gorm:"foreignKey:SourceID;references:ID"`
	Key01    string   `json:"key01" gorm:"column:col_key01;comment:预留字段"`
	Key02    string   `json:"key02" gorm:"column:col_key02;comment:预留字段"`
	Key03    string   `json:"key03" gorm:"column:col_key03;comment:预留字段"`
	Key04    string   `json:"key04" gorm:"column:col_key04;comment:预留字段"`
	Key05    string   `json:"key05" gorm:"column:col_key05;comment:预留字段"`
	Key06    string   `json:"key06" gorm:"column:col_key06;comment:预留字段"`
	Key07    string   `json:"key07" gorm:"column:col_key07;comment:预留字段"`
	Key08    string   `json:"key08" gorm:"column:col_key08;comment:预留字段"`
	Key09    string   `json:"key09" gorm:"column:col_key09;comment:预留字段"`
	Key10    string   `json:"key10" gorm:"column:col_key10;comment:预留字段"`
	Key11    string   `json:"key11" gorm:"column:col_key11;comment:预留字段"`
	Key12    string   `json:"key12" gorm:"column:col_key12;comment:预留字段"`
	Key13    string   `json:"key13" gorm:"column:col_key13;comment:预留字段"`
	Key14    string   `json:"key14" gorm:"column:col_key14;comment:预留字段"`
	Key15    string   `json:"key15" gorm:"column:col_key15;comment:预留字段"`
	Key16    string   `json:"key16" gorm:"column:col_key16;comment:预留字段"`
	Key17    string   `json:"key17" gorm:"column:col_key17;comment:预留字段"`
	Key18    string   `json:"key18" gorm:"column:col_key18;comment:预留字段"`
	Key19    string   `json:"key19" gorm:"column:col_key19;comment:预留字段"`
	Key20    string   `json:"key20" gorm:"column:col_key20;comment:预留字段"`
	Key21    string   `json:"key21" gorm:"column:col_key21;comment:预留字段"`
	Key22    string   `json:"key22" gorm:"column:col_key22;comment:预留字段"`
	Key23    string   `json:"key23" gorm:"column:col_key23;comment:预留字段"`
	Key24    string   `json:"key24" gorm:"column:col_key24;comment:预留字段"`
	Key25    string   `json:"key25" gorm:"column:col_key25;comment:预留字段"`
	Key26    string   `json:"key26" gorm:"column:col_key26;comment:预留字段"`
	Key27    string   `json:"key27" gorm:"column:col_key27;comment:预留字段"`
	Key28    string   `json:"key28" gorm:"column:col_key28;comment:预留字段"`
	Key29    string   `json:"key29" gorm:"column:col_key29;comment:预留字段"`
	Key30    string   `json:"key30" gorm:"column:col_key30;comment:预留字段"`
}

type MRules struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MRule
}

func (*MRule) TableName() string {
	return "tb_rules"
}

func (rs *MRules) Fetch() error {
	if rs.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则失败", "models:rule.Fetch()")
		return err
	}
	err := rs.PQ.Query(rs.TX, rs.All).Error
	base.NewLog("", err, "获取规则", "models:rule.Fetch()")
	return err
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

	if r.Operator != RuleOpEqual && r.Operator != RuleOpGreatThan && r.Operator != RuleOpInclude && r.Operator != RuleOpLessThan && r.Operator != RuleOpRegex {
		err := errors.New("wrong rule compare operator")
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

	err := r.TX.Create(r).Error
	base.NewLog("", err, "新增规则", "models:rule.Add()")
	return err
}
