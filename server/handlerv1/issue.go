package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/duty"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddIssueHandling(ctx *gin.Context) {
	var ih = models.MIssueHandling{}
	if err := ctx.BindJSON(&ih); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	ih.TX = base.DB().Begin()
	if err := duty.AddIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddIssueHandling, nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, ih, nil))
}

func DeleteIssueHandling(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ih = models.MIssueHandling{ID: uint(id), TX: base.DB().Begin()}
	if err := duty.DeleteIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToDeleteIssueHandling, nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func CloseIssueHandling(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ih = models.MIssueHandling{ID: uint(id), TX: base.DB().Begin()}
	if err := duty.CloseIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToCloseIssueHandling, nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func UpdateIssueHandling(ctx *gin.Context) {
	var ih = models.MIssueHandling{}
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.BindJSON(&ih); err != nil || id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	ih.ID = uint(id)
	ih.TX = base.DB().Begin()
	if err := duty.UpdateIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToUpdateIssueHandling, nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func GetIssueHandling(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ih = models.MIssueHandling{ID: uint(id), TX: base.DB().Begin()}
	if err := duty.GetIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetIssueHandling, nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, ih, nil))
}

func FetchIssueHandlings(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ihs = models.MIssueHandlings{TX: base.DB().Begin()}
	if err := duty.FetchIssueHandlings(&ihs); err != nil {
		ihs.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchIssueHandling, nil, nil))
		return
	}
	ihs.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, ihs.All, map[string]interface{}{"total": pq.Total}))
}

func FetchIssueHandlingEvents(ctx *gin.Context) {
	eventid := ctx.Param("eventid")
	if len(eventid) == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var ih = models.MIssueHandling{EventId: eventid, TX: base.DB().Begin()}
	if err := duty.FetchIssueHandlingEvents(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchIssueHandlingEvents, nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, ih.Events, nil))
}
