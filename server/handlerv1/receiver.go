package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/receiver"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddReceiver(ctx *gin.Context) {
	var rcv = models.MReceiver{}
	if err := ctx.BindJSON(&rcv); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	rcv.TX = base.DB()
	if err := receiver.Add(&rcv); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddReceiver, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, rcv, nil))
}

func FetchReceivers(ctx *gin.Context) {
	var rcvs = models.MReceivers{}
	if err := ctx.BindQuery(&rcvs.PQ); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	rcvs.TX = base.DB()
	if err := receiver.Fetch(&rcvs); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchReceivers, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, rcvs.All, map[string]interface{}{"total": rcvs.PQ.Total}))
}

func GetReceiver(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var rcv = models.MReceiver{ID: uint(id), TX: base.DB()}
	if err := receiver.Get(&rcv); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetReceiver, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, rcv, nil))
}

func DeleteReceiver(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var rcv = models.MReceiver{ID: uint(id), TX: base.DB().Begin()}
	if err := rcv.Delete(); err != nil {
		rcv.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToDeleteReceiver, nil, nil))
		return
	}
	rcv.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}
