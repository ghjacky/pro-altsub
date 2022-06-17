package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/source"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddSource(ctx *gin.Context) {
	var src = &models.MSource{}
	if err := ctx.Bind(src); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	if len(src.Name) <= 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
		return
	}
	src.BaseModel.TX = base.DB()
	if err := source.Add(src); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, src, nil))
}

func FetchSources(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	if srcs, err := source.Fetch(base.DB(), &pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToQuerySources, nil, nil))
	} else {
		ctx.JSON(http.StatusOK, newHttpResponse(nil, srcs.All, map[string]interface{}{"total": srcs.PQ.Total}))
	}
}
