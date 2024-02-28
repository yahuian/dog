package errcode

import (
	"errors"
	"net/http"
)

type errCode struct {
	code int // response body code
	err  error
}

func (e *errCode) Error() string {
	return e.err.Error()
}

func newErrcode(code int, msg string) error {
	return &errCode{
		code: code,
		err:  errors.New(msg),
	}
}

func Parse(e error) (int, string) {
	for {
		if e == nil {
			return 0, ""
		}

		err, ok := e.(*errCode)
		if !ok {
			e = errors.Unwrap(e)
			continue
		}

		return err.code, err.Error()
	}
}

// Code2Status 业务码到 http 状态码
func Code2Status(c int) int {
	// 业务码
	if c >= 40000 && c < 50000 {
		return http.StatusBadRequest
	}
	if c > 50000 {
		return http.StatusInternalServerError
	}

	// 复用状态码
	return c
}

// BadRequest 通用的客户端请求错误，一般是参数校验失败
// 正常情况由于前端也做了校验，一般来说错误信息用户看不到，所以报错直接用英文
func BadRequest(err error) error {
	return &errCode{
		code: http.StatusBadRequest,
		err:  err,
	}
}

// Server 服务端内部错误不细分，记录清楚日志方便排查即可
func Server(err error) error {
	return &errCode{
		code: http.StatusInternalServerError,
		err:  err,
	}
}

var (
	// 401 403 404 一般前端会有独立的页面，所以也用英文
	UnauthorizedErr = newErrcode(http.StatusUnauthorized, "unauthorized")
	ForbiddenErr    = newErrcode(http.StatusForbidden, "forbidden")
	NotFoundErr     = newErrcode(http.StatusNotFound, "not found")

	// 业务细分错误用中文，方便直接返回给前端提醒用户
	// user 401xx
)
