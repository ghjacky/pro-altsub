package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/duty"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchDuties(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "FetchDuties()", nil, nil))
		return
	}
	var ds = models.MDuties{}
	ds.TX = base.DB().Begin()
	ds.PQ = pq
	if err := duty.Fetch(&ds); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchDuties, "FetchDuties()", nil, nil))
		ds.TX.Rollback()
		return
	}
	ds.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "FetchDuties()", ds.All, map[string]interface{}{"total": ds.PQ.Total}))
}

func AddDuty(ctx *gin.Context) {
	var d = models.MDuty{}
	if err := ctx.BindJSON(&d); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "AddDuty()", nil, nil))
		return
	}
	d.TX = base.DB().Begin()
	if err := duty.Add(&d); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddDuty, "AddDuty()", nil, nil))
		d.TX.Rollback()
		return
	}
	d.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "AddDuty()", nil, nil))
}
