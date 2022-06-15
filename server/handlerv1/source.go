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
		ctx.JSON(http.StatusOK, NewHttpResponse(models.ErrorCodeBadRequest, nil, nil, nil))
		return
	}
	if len(src.Name) <= 0 {
		ctx.JSON(http.StatusOK, NewHttpResponse(models.ErrorCodeEmptySource, nil, nil, nil))
		return
	}
	src.BaseModel.DB = base.DB()
	if err := source.Add(src); err != nil {
		ctx.JSON(http.StatusOK, NewHttpResponse(models.ErrorCodeEmptySource, err, nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, NewHttpResponse(0, nil, src, nil))
}
