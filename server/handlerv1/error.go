package handlerv1

type ResponseError struct {
	code    int
	message string
	error
}

const (
	ErrCodeInternal = -1
)

var (
	ErrorOther ResponseError = ResponseError{
		code:    9999,
		message: "未知错误",
	}
	ErrorBadRequest ResponseError = ResponseError{
		code:    9994,
		message: "参数错误",
	}
	ErrorEmptySource ResponseError = ResponseError{
		code:    1001,
		message: "source为空",
	}
	ErrorFailedToAddSource ResponseError = ResponseError{
		code:    1002,
		message: "告警源接入失败",
	}
	ErrorFailedToQuerySources ResponseError = ResponseError{
		code:    1003,
		message: "告警源查询失败",
	}
	ErrorEmptySchemaData ResponseError = ResponseError{
		code:    2001,
		message: "schema数据为空",
	}
	ErrorFailedToAddSchema ResponseError = ResponseError{
		code:    2002,
		message: "新增schema失败",
	}
	ErrorFailedToWriteEvent ResponseError = ResponseError{
		code:    3002,
		message: "事件写入失败",
	}
)

func NewErr(code int, message string, e error) *ResponseError {
	return &ResponseError{
		code:    code,
		message: message,
		error:   e,
	}
}

func (err *ResponseError) Code() int {
	return err.code
}

func (err *ResponseError) Message() string {
	return err.message
}

func (err *ResponseError) Error() string {
	if err.error != nil {
		return err.error.Error()
	} else {
		return ""
	}
}
