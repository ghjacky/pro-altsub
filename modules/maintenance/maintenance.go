package maintenance

import (
	"altsub/base"
	"altsub/models"
	"errors"
	"time"
)

func Check(rs []models.MRule) bool {
	for _, r := range rs {
		r.TX = base.DB()
		if err := r.GetChain("Prev", "Clauses", "Maintenances"); err != nil {
			base.NewLog("error", err, "维护项检查出错", "maintenance:Check()")
			return false
		}
		var _r = r
		for {
			if len(_r.Maintenances) > 0 {
				// 检测是否过期
				for _, _m := range _r.Maintenances {
					now := time.Now().Local().Unix()
					if _m.StartAt > now || _m.EndAt < now {
						base.NewLog("info", nil, "根据规则检测到相关维护项，但不在生效时间范围内", "maintenance:Check()")
					} else {
						base.NewLog("info", nil, "根据规则检测到相关维护项", "maintenance:Check()")
						return true
					}
				}
			}
			if _r.Prev != nil {
				_r = *_r.Prev
			} else {
				break
			}
		}
	}
	return false
}

func Fetch(ms *models.MMaintenances) error {
	if ms.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取维护项列表失败", "maintenance:Fetch()")
		return err
	}
	return ms.Fetch()
}

func Add(m *models.MMaintenance) error {
	if m.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增维护项失败", "maintenance:Add()")
		return err
	}
	return m.Add()
}

func Remove(m *models.MMaintenance) error {
	if m.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "删除维护项失败", "maintenance:Remove()")
		return err
	}
	return m.Delete()
}

func Get(m *models.MMaintenance) error {
	if m.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取维护项失败", "maintenance:Get()")
		return err
	}
	return m.Get()
}
