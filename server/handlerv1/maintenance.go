package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/maintenance"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchMaintenances(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ms = models.MMaintenances{TX: base.DB(), PQ: pq}
	if err := maintenance.Fetch(&ms); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchMaintenances, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, ms.All, map[string]interface{}{"total": ms.PQ.Total}))
}

func AddMaintenance(ctx *gin.Context) {
	var m = models.MMaintenance{}
	if err := ctx.BindJSON(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	m.TX = base.DB()
	if err := maintenance.Add(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddMaintenance, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, m, nil))
}

func RemoveMaintenance(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var m = models.MMaintenance{TX: base.DB()}
	m.ID = uint(id)
	if err := maintenance.Remove(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToRemoveMaintenance, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func GetMaintenance(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var m = models.MMaintenance{TX: base.DB()}
	m.ID = uint(id)
	if err := maintenance.Get(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetMaintenance, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, m, nil))
}
