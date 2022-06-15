package models

type ErrorCode int

const (
	ErrorCodeUnkown     ErrorCode = 9999
	ErrorCodeBadRequest ErrorCode = 9994

	ErrorCodeEmptySource       ErrorCode = 1001
	ErrorCodeFailedToAddSource ErrorCode = 1002
)

func (err ErrorCode) String() string {
	switch err {
	case ErrorCodeBadRequest:
		return "参数错误"
	case ErrorCodeFailedToAddSource:
		return "告警源接入失败"
	case ErrorCodeEmptySource:
		return "source为空"
	default:
		return "未知错误"
	}
}

func (err ErrorCode) Code() int {
	return int(err)
}
