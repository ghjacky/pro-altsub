package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/event"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveRawEvent(ctx *gin.Context) {
	var  srcName = ctx.Query("source")
	if len( srcName) <= 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
		return
	}
	var ev = models.MEvent{}
	if err := ctx.Bind(&ev.Data); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	base.NewLog("trace", nil, fmt.Sprintf("接收到原始告警消息: %s", string(ev.Data)), "handlerv1:ReceiveRawEvent()")
	if err := event.Receive( srcName, ev); err != nil {
		base.NewLog("error", err, "事件写入失败", "handlerv1:ReceiveRawEvent()")
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToWriteEvent, nil, nil))
		//
		// TODO：kafka 数据写失败 告警
		//

		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}
