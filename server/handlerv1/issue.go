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
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "AddIssueHandling()", nil, nil))
		return
	}
	ih.TX = base.DB().Begin()
	if err := duty.AddIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddIssueHandling, "AddIssueHandling()", nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "AddIssueHandling()", ih, nil))
}

func DeleteIssueHandling(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "DeleteIssueHandling()", nil, nil))
		return
	}
	var ih = models.MIssueHandling{ID: uint(id), TX: base.DB().Begin()}
	if err := duty.DeleteIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToDeleteIssueHandling, "DeleteIssueHandling()", nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "DeleteIssueHandling()", nil, nil))
}

func CloseIssueHandling(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "CloseIssueHandling()", nil, nil))
		return
	}
	var ih = models.MIssueHandling{ID: uint(id), TX: base.DB().Begin()}
	if err := duty.CloseIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToCloseIssueHandling, "CloseIssueHandling()", nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "CloseIssueHandling()", nil, nil))
}

func UpdateIssueHandling(ctx *gin.Context) {
	var ih = models.MIssueHandling{}
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.BindJSON(&ih); err != nil || id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "UpdateIssueHandling()", nil, nil))
		return
	}
	ih.ID = uint(id)
	ih.TX = base.DB().Begin()
	if err := duty.UpdateIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToUpdateIssueHandling, "UpdateIssueHandling()", nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "UpdateIssueHandling()", nil, nil))
}

func GetIssueHandling(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "GetIssueHandling()", nil, nil))
		return
	}
	var ih = models.MIssueHandling{ID: uint(id), TX: base.DB().Begin()}
	if err := duty.GetIssueHandling(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetIssueHandling, "GetIssueHandling()", nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "GetIssueHandling()", ih, nil))
}

func FetchIssueHandlings(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "FetchIssueHandlings()", nil, nil))
		return
	}
	var ihs = models.MIssueHandlings{TX: base.DB().Begin()}
	if err := duty.FetchIssueHandlings(&ihs); err != nil {
		ihs.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchIssueHandling, "FetchIssueHandlings()", nil, nil))
		return
	}
	ihs.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "FetchIssueHandlings()", ihs.All, map[string]interface{}{"total": pq.Total}))
}

func FetchIssueHandlingEvents(ctx *gin.Context) {
	eventid := ctx.Param("eventid")
	if len(eventid) == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "FetchIssueHandlingEvents()", nil, nil))
		return
	}
	var ih = models.MIssueHandling{EventId: eventid, TX: base.DB().Begin()}
	if err := duty.FetchIssueHandlingEvents(&ih); err != nil {
		ih.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchIssueHandlingEvents, "FetchIssueHandlingEvents()", nil, nil))
		return
	}
	ih.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "FetchIssueHandlingEvents()", ih.Events, nil))
}
