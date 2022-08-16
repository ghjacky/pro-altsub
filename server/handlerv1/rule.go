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
	var rls = []models.MRule{}
	name := ctx.Query("rule_name")
	if len(name) <= 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptyRuleName, nil, nil))
		return
	}
	if err := ctx.Bind(&rls); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	db := base.DB().Begin()
	for _, rl := range rls {
		rl.Name = name
		if len(rl.Source.Name) <= 0 {
			ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
			return
		}
		rl.Source.TX = db
		rl.TX = rl.Source.TX
		if err := rule.Add(&rl); err != nil {
			db.Rollback()
			ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddRule, nil, nil))
			return
		}
	}
	db.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
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

func FetchRuleChain(ctx *gin.Context) {
	var r = models.MRule{}
	if err := ctx.BindJSON(&r); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	if r.ID == 0 && (len(r.Name) <= 0 || r.SourceID == 0) {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptyRuleNameOrZeroSourceID, nil, nil))
		return
	}
	r.TX = base.DB()
	if err := rule.FetchRuleChain(&r); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToFetchRuleChain, nil, nil))
	} else {
		ctx.JSON(http.StatusOK, newHttpResponse(nil, r, nil))
	}
}

func GetRule(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var r = models.MRule{ID: uint(id), TX: base.DB()}
	if err := rule.Get(&r); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetRule, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, r, nil))
}

func DeleteRule(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var r = models.MRule{ID: uint(id), TX: base.DB()}
	if err := rule.Delete(&r); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToDeleteRule, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}
