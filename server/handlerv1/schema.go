package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddSchema(ctx *gin.Context) {
	var scm = models.MSchema{}
	sourceName := ctx.Query("source_name")
	sourceId, _ := strconv.Atoi(ctx.Query("source_id"))
	if len(sourceName) <= 0 && sourceId == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
		return
	}
	if err := ctx.Bind(&scm); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	src := models.MSource{Name: sourceName}
	src.ID = uint(sourceId)
	src.TX = base.DB()
	scm.Source = src
	scm.TX = src.TX
	if err := schema.Add(&scm); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddSchema, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, scm, nil))
}

func UpdateSchema(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var schm = models.MSchema{}
	if err := ctx.BindJSON(&schm);id == 0 || err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	schm.ID = uint(id)
	schm.TX = base.DB().Begin()
	if err := schema.Update(&schm); err != nil {
		schm.TX.Rollback()
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToUpdateSchema, nil, nil))
		return
	}
	schm.TX.Commit()
	ctx.JSON(http.StatusOK, newHttpResponse(nil, nil, nil))
}

func FetchSchemas(ctx *gin.Context) {
	var pq = models.PageQuery{}
	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	if schms, err := schema.Fetch(base.DB(), &pq); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToQuerySchemas, nil, nil))
	} else {
		ctx.JSON(http.StatusOK, newHttpResponse(nil, schms.All, map[string]interface{}{"total": schms.PQ.Total}))
	}
}

func GetSchema(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id == 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	var schm = models.MSchema{ID: uint(id), TX: base.DB()}
	if err := schema.Get(&schm); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToGetSchema, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, schm, nil))
}
