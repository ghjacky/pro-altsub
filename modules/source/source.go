package source

import (
	"altsub/base"
	"altsub/models"

	"gorm.io/gorm"
)

func Fetch(tx *gorm.DB, pq *models.PageQuery) (*models.MSources, error) {
	if tx == nil {
		tx = base.DB()
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
	return ss, ss.Fetch()
}

func Add(source *models.MSource) error {
	return source.Add()
}
