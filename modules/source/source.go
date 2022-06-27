package source

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/event"
	"errors"

	"gorm.io/gorm"
)

func Fetch(tx *gorm.DB, pq *models.PageQuery) (*models.MSources, error) {
	if tx == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取告警源失败", "source:Fetch()")
		return nil, err
	}
	if pq == nil {
		pq = &models.PageQuery{}
		pq.Page = 1
		pq.Size = 10000
		pq.Order = "+id"
	}
	var ss = &models.MSources{
		TX:  tx,
		PQ:  *pq,
		All: []*models.MSource{},
	}
	return ss, ss.Fetch("Rules")
}

func Add(src *models.MSource) (err error) {
	if src.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "新增数据源失败", "source:Add()")
		return err
	}
	if err := src.Add(); err != nil {
		return err
	} else {
		// 新增source，需要后台启动从对应topic消费事件
		_ = base.ReadFromKafka(src.Name)
		// 并从该source对应的message buffer中消费并处理事件
		event.ReadAndParseEventFromBufferForever(src.Name)
		return nil
	}
}

func GetByName(src *models.MSource) (err error) {
	if src.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "根据名称获取数据源失败", "source:GetByName()")
		return err
	}
	if len(src.Name) <= 0 {
		err := errors.New("empty source name")
		base.NewLog("error", err, "根据名称获取数据源失败", "source:GetByName()")
		return err
	}
	return src.GetByName()
}

func Get(src *models.MSource) error {
	if src.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取数据源失败", "source:Get()")
		return err
	}
	return src.Get()
}