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
	ErrorFailedToGetSource ResponseError = ResponseError{
		code:    1004,
		message: "告警源查询失败",
	}
	ErrorFailedToFetchSourceTypes ResponseError = ResponseError{
		code:    1005,
		message: "获取告警源类型失败",
	}
	ErrorEmptySchemaData ResponseError = ResponseError{
		code:    2001,
		message: "schema数据为空",
	}
	ErrorFailedToAddSchema ResponseError = ResponseError{
		code:    2002,
		message: "新增schema失败",
	}
	ErrorFailedToQuerySchemas ResponseError = ResponseError{
		code:    2003,
		message: "schema查询失败",
	}
	ErrorFailedToGetSchema ResponseError = ResponseError{
		code:    2004,
		message: "schema查询失败",
	}
	ErrorFailedToUpdateSchema ResponseError = ResponseError{
		code:    2005,
		message: "schema更新失败",
	}
	ErrorFailedToWriteEvent ResponseError = ResponseError{
		code:    3002,
		message: "事件写入失败",
	}
	ErrorEmptyRuleName ResponseError = ResponseError{
		code:    4002,
		message: "规则名称为空",
	}
	ErrorEmptyRuleNameOrZeroSourceID ResponseError = ResponseError{
		code:    4001,
		message: "规则名称为空或sourceId为0",
	}
	ErrorFailedToFetchRuleChain ResponseError = ResponseError{
		code:    4003,
		message: "获取规则链失败",
	}
	ErrorFailedToGetRule ResponseError = ResponseError{
		code:    4004,
		message: "获取规则详情失败",
	}
	ErrorFailedToAddRule ResponseError = ResponseError{
		code:    4005,
		message: "新增规则失败",
	}
	ErrorFailedToDeleteRule ResponseError = ResponseError{
		code:    4006,
		message: "规则删除失败",
	}
	ErrorFailedToAddReceiver ResponseError = ResponseError{
		code:    5002,
		message: "新增接收者失败",
	}
	ErrorFailedToFetchReceivers ResponseError = ResponseError{
		code:    5003,
		message: "获取接收者列表失败",
	}
	ErrorFailedToGetReceiver ResponseError = ResponseError{
		code:    5004,
		message: "获取接收者详情失败",
	}
	ErrorFailedToDeleteReceiver ResponseError = ResponseError{
		code:    5005,
		message: "删除接收者失败",
	}
	ErrorFailedToSubscribeRules ResponseError = ResponseError{
		code:    6002,
		message: "订阅规则失败",
	}
	ErrorFailedToAssignRules ResponseError = ResponseError{
		code:    6003,
		message: "指派规则失败",
	}
	ErrorFailedToFetchSubscribes ResponseError = ResponseError{
		code:    6004,
		message: "获取订阅/指派关系失败",
	}
	ErrorFailedToAddDuty ResponseError = ResponseError{
		code:    7002,
		message: "新增排班班次失败",
	}
	ErrorFailedToFetchDuties ResponseError = ResponseError{
		code:    7003,
		message: "获取排班班次列表失败",
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
