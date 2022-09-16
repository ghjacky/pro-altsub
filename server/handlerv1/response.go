package handlerv1

import (
	"altsub/base"
)

type HttpResponse map[string]interface{}

func newHttpResponse(err *ResponseError, caller string, data interface{}, extras map[string]interface{}) HttpResponse {
	var hr = HttpResponse{}
	if err == nil {
		hr["code"] = 0
		hr["message"] = "ok"
	} else {
		hr["code"] = err.Code()
		hr["message"] = err.Message()
	}
	hr["data"] = data
	for k, v := range extras {
		hr[k] = v
	}
	msg, _ := hr["message"].(string)
	if err != nil {
		base.NewLog("error", err, msg, caller)
	}
	return hr
}
