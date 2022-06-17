package rule

import (
	"altsub/base"
	"altsub/models"
	"errors"

	"gorm.io/gorm"
)

func Add(rl *models.MRule) (err error) {
	if rl == nil {
		err = errors.New("empty rule data")
		base.NewLog("error", err, "新增规则失败", "rule:Add()")
		return
	}
	if rl.BaseModel.TX == nil {
		err = errors.New("nil db object")
		base.NewLog("error", err, "新增规则失败", "rule:Add()")
		return
	}
	return rl.Add()
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
		All: []*models.MRule{},
	}
	return rs, rs.Fetch()
}