package base

import (
	"errors"
	"strconv"
)

func ParseFloat64(v interface{}) (float64, error) {
	switch _v := v.(type) {
	case float32:
		return float64(_v), nil
	case float64:
		return _v, nil
	case int:
		return float64(_v), nil
	case int64:
		return float64(_v), nil
	case int32:
		return float64(_v), nil
	case string:
		f, e := strconv.ParseFloat(_v, 64)
		if e != nil {
			return 0, e
		} else {
			return f, nil
		}
	default:
		return 0, errors.New("unknown type of value")
	}
}
