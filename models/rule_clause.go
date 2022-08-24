package models

import (
	"errors"

	"gorm.io/gorm"
)

const (
	RuleOpEqual = iota + 1
	RuleOpGreatThan
	RuleOpLessThan
	RuleOpInclude
	RuleOpRegex
)

type MRuleClause struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Key      string `json:"key" gorm:"column:col_key;type:varchar(32);not null;comment:规则项比对key"`
	Operator int    `json:"operator" gorm:"column:col_operator;not null;default:1;comment:规则项比对操作符"`
	Value    string `json:"value" gorm:"column:col_value;type:varchar(64);not null;comment:规则项比对value"`
	RuleID   uint   `json:"rule_id" gorm:"column:col_rule_id;not null"`
}

type MRuleClauses struct {
	PQ  PageQuery
	TX  *gorm.DB
	All []*MRuleClause
}

func (*MRuleClause) TableName() string {
	return "tb_rule_clauses"
}

func (rc *MRuleClause) BeforeCreate(tx *gorm.DB) error {
	if rc.Operator <= 0 || rc.Operator > RuleOpRegex {
		return errors.New("wrong rule clause operator")
	}
	return nil
}
