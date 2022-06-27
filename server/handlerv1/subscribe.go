package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/subscribe"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Subscribe(ctx *gin.Context) {
	rcvId, _ := strconv.Atoi(ctx.Param("id"))
	if rcvId == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var sub = models.MSubscribe{}
	if err := ctx.BindQuery(&sub); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	sub.Type = models.SubscribeTypeSub
	var rs = []models.MRule{}
	if err := ctx.BindJSON(&rs); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var rcv = models.MReceiver{}
	rcv.TX = base.DB().WithContext(context.WithValue(context.Background(), "subscribe", sub))
	rcv.ID = uint(rcvId)
	if err := subscribe.Subscribe(rcv, rs); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToSubscribeRules, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func Assign(ctx *gin.Context) {
	rId, _ := strconv.Atoi(ctx.Param("id"))
	if rId == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var sub = models.MSubscribe{}
	if err := ctx.BindQuery(&sub); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	sub.Type = models.SubscribeTypeAss
	var rcvs = []models.MReceiver{}
	if err := ctx.BindJSON(&rcvs); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var r = models.MRule{}
	r.TX = base.DB().WithContext(context.WithValue(context.Background(), "subscribe", sub))
	r.ID = uint(rId)
	if err := subscribe.Assign(r, rcvs); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAssignRules, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func FetchSubscribe(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ss = models.MSubscribes{TX: base.DB().Begin(), PQ: pq}
	if err := subscribe.Fetch(&ss); err != nil {
		ss.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchSubscribes, nil, nil))
		return
	}
	ss.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, ss.All, map[string]interface{}{"total": ss.PQ.Total}))
}
