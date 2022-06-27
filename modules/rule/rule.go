package rule

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

func Delete(rl *models.MRule) error {
	if rl.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "规则删除失败", "rule:Delete()")
		return err
	}
	return rl.Delete()
}

func Add(rl *models.MRule) error {
	if rl == nil {
		err := errors.New("empty rule data")
		base.NewLog("error", err, "新增规则失败", "rule:Add()")
		return err
	}
	if rl.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增规则失败", "rule:Add()")
		return err
	}
	if err := rl.Add(); err != nil {
		base.NewLog("error", err, "新增规则失败", "rule:Add()")
		return err
	}
	go func() {
		// 新增规则需要对已有事件进行规则检测追加（默认应用最近一个月的事件）
		if err := checkAndAppend(rl); err != nil {
			base.NewLog("warn", err, "检查并追加规则失败", "rule:Add()")
		}
	}()
	return nil
}

func Fetch(tx *gorm.DB, pq *models.PageQuery) (*models.MRules, error) {
	if tx == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则失败", "rule:Fetch()")
		return nil, err
	}
	if pq == nil {
		pq = &models.PageQuery{}
		pq.Page = 1
		pq.Size = 10000
		pq.Order = "+id"
	}
	var rs = &models.MRules{
		TX:  tx,
		PQ:  *pq,
		All: []models.MRule{},
	}
	return rs, rs.Fetch()
}

func FetchRuleChain(r *models.MRule) error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则链失败", "rule:FetchRuleChain()")
		return err
	}
	if r.ID == 0 && (len(r.Name) <= 0 || r.SourceID == 0) {
		err := errors.New("empty rule id and (empty name or source id)")
		base.NewLog("error", err, "获取规则链失败", "rule:FetchRuleChain()")
		return err
	}
	return r.GetChain("Prev", "Clauses")
}

func Get(r *models.MRule) error {
	if r.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取规则想起失败", "rule:Get()")
		return err
	}
	return r.Get("Receivers")
}

// 检查并追加规则（新增规则后）
const (
	defaultCheckDurationLatest = "360h"
)

func checkAndAppend(r *models.MRule) error {
	dur, e := time.ParseDuration(defaultCheckDurationLatest)
	if e != nil {
		base.NewLog("error", e, "检查并追加规则失败", "event:CheckAndAppendRule()")
		return e
	}
	t := time.Now().Local().Add(-1 * dur)
	var evs = models.MSchemaedEvents{TX: base.DB().Begin()}
	if err := evs.FetchAfter(t); err != nil {
		base.NewLog("error", err, "检查并追加规则失败", "event:CheckAndAppendRule()")
		evs.TX.Rollback()
		return err
	}
	for _, ev := range evs.All {
		ev.TX = evs.TX
		var schmEv = schema.SchemaedEvent{}
		if err := json.Unmarshal(ev.Data, &schmEv); err != nil {
			base.NewLog("error", err, "检查并追加规则失败", "event:CheckAndAppendRule()")
			ev.TX.Rollback()
			return err
		}
		rs := Check(r.Source.Name, schmEv)
		for _, r := range rs {
			if err := ev.AppendRule(r); err != nil {
				base.NewLog("error", err, "检查并追加规则失败", "event:CheckAndAppendRule()")
				ev.TX.Rollback()
				return err
			}
		}
	}
	evs.TX.Commit()
	return nil
}

func Check(srcName string, ev schema.SchemaedEvent) []models.MRule {
	var rs = models.MRules{}
	var src = models.MSource{}
	src.TX = base.DB()
	src.Name = srcName
	if err := src.GetByName(); err != nil {
		base.NewLog("error", err, "检测事件相关维护项失败", "event:CheckMaintenance()")
		return nil
	}
	rs.PQ = models.PageQuery{
		Size:   9999,
		Page:   1,
		Search: fmt.Sprintf("col_source_id:%d", src.ID),
	}
	rs.TX = src.TX
	if err := rs.Fetch("Receivers", "Clauses"); err != nil {
		base.NewLog("error", err, "检测事件相关维护项失败", "event:CheckMaintenance()")
		return nil
	} else {
		var all = []models.MRule{}
		// 根据key check 事件是否符合规则
		m := map[string]interface{}{}
		for _, k := range ev {
			m[k.Key] = k.Value
		}
		for _, r := range rs.All {
			rmatch := false
			for i, rc := range r.Clauses {
				rcmatch := true
				if v, ok := m[rc.Key]; !ok {
					rcmatch = false
				} else {
					// check kv ok
					switch rc.Operator {
					case models.RuleOpEqual:
						rcmatch = (v == rc.Value)
					case models.RuleOpGreatThan:
						_v, e := base.ParseFloat64(v)
						_rcv, _e := base.ParseFloat64(rc.Value)
						if e != nil {
							base.NewLog("error", e, "检测维护项失败，操作符错误", "event:CheckMaintenance()")
							rcmatch = false
						} else if _e != nil {
							base.NewLog("error", _e, "检测维护项失败，操作符错误", "event:CheckMaintenance()")
							rcmatch = false
						} else {
							rcmatch = (_v > _rcv)
						}
					case models.RuleOpLessThan:
						_v, e := base.ParseFloat64(v)
						_rcv, _e := base.ParseFloat64(rc.Value)
						if e != nil {
							base.NewLog("error", e, "检测维护项失败，操作符错误", "event:CheckMaintenance()")
							rcmatch = false
						} else if _e != nil {
							base.NewLog("error", _e, "检测维护项失败，操作符错误", "event:CheckMaintenance()")
							rcmatch = false
						} else {
							rcmatch = (_v < _rcv)
						}
					case models.RuleOpInclude:
						_v, ok := v.(string)
						if !ok {
							rcmatch = false
						} else {
							rcmatch = (strings.Contains(_v, rc.Value))
						}
					case models.RuleOpRegex:
						_v, ok := v.(string)
						if !ok {
							rcmatch = false
						} else {
							p, e := regexp.Compile(_v)
							if e != nil {
								rcmatch = false
							} else {
								rcmatch = p.Match([]byte(_v))
							}
						}
					default:
						base.NewLog("error", errors.New("wrong operator"), "检测维护项失败，操作符错误", "event:CheckMaintenance()")
						rcmatch = false
					}
				}
				switch r.Logic {
				case models.RuleLogicIntersection:
					rmatch = (rmatch || (i == 0)) && rcmatch
				case models.RuleLogicConcatenation:
					rmatch = rmatch || rcmatch
				}
			}
			if rmatch {
				all = append(all, r)
			}
		}
		return all
	}
}
