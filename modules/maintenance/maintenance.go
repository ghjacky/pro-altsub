package maintenance

import (
	"altsub/base"
	"altsub/models"
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
