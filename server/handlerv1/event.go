package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/event"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveRawEvent(ctx *gin.Context) {
	srcName, exist := ctx.GetQuery("source")
	if exist && len(srcName) <= 0 {
		base.NewLog("error", errors.New("empty source name"), "source名称为空", "ReceiveRawEvent()")
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, "ReceiveRawEvent()", nil, nil))
		return
	} else if !exist && len(srcName) <= 0 {
		base.NewLog("error", errors.New("no source"), "缺少source参数", "ReceiveRawEvent()")
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, "ReceiveRawEvent()", nil, nil))
		return
	}
	var ev = models.MEvent{}
	if err := ctx.BindJSON(&ev.Data); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "ReceiveRawEvent()", err, nil))
		return
	}
	base.NewLog("trace", nil, fmt.Sprintf("接收到事件数据：: %s", string(ev.Data)), "handlerv1:ReceiveRawEvent()")
	if err := event.Receive(srcName, ev); err != nil {
		base.NewLog("error", err, "事件写入失败", "handlerv1:ReceiveRawEvent()")
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToWriteEvent, "ReceiveRawEvent()", nil, nil))
		//
		// TODO：kafka 数据写失败 告警
		//

		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "ReceiveRawEvent()", nil, nil))
}
