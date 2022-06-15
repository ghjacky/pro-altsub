package handlerv1

import (
	"altsub/base"
	"altsub/models"
)

type HttpResponse map[string]interface{}

func NewHttpResponse(code models.ErrorCode, err error, data interface{}, extras map[string]interface{}) HttpResponse {
	var hr = HttpResponse{}
	hr["code"] = code.Code()
	hr["message"] = code.String()
	hr["data"] = data
	for k, v := range extras {
		hr[k] = v
	}
	msg, _ := hr["message"].(string)
	if err != nil {
		base.NewLog("error", err, msg, "newHttpResponse()")
	}
	return hr
}
