package handlerv1

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
	"altsub/modules/source"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddSchema(ctx *gin.Context) {
	var scm = models.MSchema{}
	sourceName := ctx.Query("source")
	if err := ctx.Bind(&scm); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorBadRequest, nil, nil))
		return
	}
	if len(sourceName) <= 0 || len(scm.Data) <= 0 {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorEmptySource, nil, nil))
		return
	}
	var src = models.MSource{Name: sourceName}
	src.BaseModel.TX = base.DB()
	if err := source.GetByName(&src); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddSchema, nil, nil))
		return
	}
	scm.Source = src
	scm.BaseModel.TX = src.BaseModel.TX
	if err := schema.Add(&scm); err != nil {
		ctx.JSON(http.StatusOK, newHttpResponse(&ErrorFailedToAddSchema, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, newHttpResponse(nil, scm, nil))
}
