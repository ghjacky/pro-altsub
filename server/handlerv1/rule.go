package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/rule"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddRule(ctx *gin.Context) {
	var rl = models.MRule{}
	sourceName := ctx.Query("source_name")
	sourceId, _ := strconv.Atoi(ctx.Query("source_id"))
	if len(sourceName) <= 0 && sourceId == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
		return
	}
	if err := ctx.Bind(&rl); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	src := &models.MSource{Name: sourceName}
	src.ID = uint(sourceId)
	src.TX = base.DB()
	rl.Source = src
	rl.BaseModel.TX = src.TX
	if err := rule.Add(&rl); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddSchema, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, rl, nil))
}

func FetchRules(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	if rls, err := rule.Fetch(base.DB(), &pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToQuerySources, nil, nil))
	} else {
		ctx.JSON(http.StatusOK, newHttpResponse(nil, rls.All, map[string]interface{}{"total": rls.PQ.Total}))
	}
}
