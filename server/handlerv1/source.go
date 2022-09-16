package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/source"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddSource(ctx *gin.Context) {
	var src = &models.MSource{}
	if err := ctx.Bind(src); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "AddSource()", nil, nil))
		return
	}
	if len(src.Name) <= 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, "AddSource()", nil, nil))
		return
	}
	src.TX = base.DB()
	if err := source.Add(src); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, "AddSource()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "AddSource()", nil, nil))
}

func FetchSources(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "FetchSources()", nil, nil))
		return
	}
	if srcs, err := source.Fetch(base.DB(), &pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToQuerySources, "FetchSources()", nil, nil))
	} else {
		ctx.JSON(http.StatusOK, newHttpResponse(nil, "FetchSources()", srcs.All, map[string]interface{}{"total": srcs.PQ.Total}))
	}
}

func GetSource(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, "GetSource()", nil, nil))
		return
	}
	var src = models.MSource{ID: uint(id), TX: base.DB()}
	if err := source.Get(&src); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetSource, "GetSource()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "GetSource()", src, nil))
}

func FetchSourceTypes(ctx *gin.Context) {
	var ss = models.MSources{TX: base.DB()}
	if err := ss.FetchTypes(); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchSourceTypes, "FetchSourceTypes()", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, "FetchSourceTypes()", ss, nil))
}
