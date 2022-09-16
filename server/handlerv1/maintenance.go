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
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "FetchMaintenances()", nil, nil))
		return
	}
	var ms = models.MMaintenances{TX: base.DB(), PQ: pq}
	if err := maintenance.Fetch(&ms); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchMaintenances, "FetchMaintenances()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "FetchMaintenances()", ms.All, map[string]interface{}{"total": ms.PQ.Total}))
}

func AddMaintenance(ctx *gin.Context) {
	var m = models.MMaintenance{}
	if err := ctx.BindJSON(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "AddMaintenance()", nil, nil))
		return
	}
	m.TX = base.DB()
	if err := maintenance.Add(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddMaintenance, "AddMaintenance()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "AddMaintenance()", m, nil))
}

func RemoveMaintenance(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "RemoveMaintenance()", nil, nil))
		return
	}
	var m = models.MMaintenance{TX: base.DB()}
	m.ID = uint(id)
	if err := maintenance.Remove(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToRemoveMaintenance, "RemoveMaintenance()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "RemoveMaintenance()", nil, nil))
}

func GetMaintenance(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "GetMaintenance()", nil, nil))
		return
	}
	var m = models.MMaintenance{TX: base.DB()}
	m.ID = uint(id)
	if err := maintenance.Get(&m); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetMaintenance, "GetMaintenance()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "GetMaintenance()", m, nil))
}
