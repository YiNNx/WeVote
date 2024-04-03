package errors

import "github.com/YiNNx/WeVote/internal/config"

type customError struct {
	Msg string
}

func newCustomError(msg string) customError {
	return customError{
		Msg: msg,
	}
}

func (e *customError) WithErrDetail(err error) *customError {
	if !config.C.Server.DebugMode {
		return e
	}
	return &customError{
		Msg: e.Msg + ": " + e.Error(),
	}
}

func (e customError) Error() string {
	return e.Msg
}

var (
	ErrInvalidParams          = newCustomError("参数错误")
	ErrInvalidUsernameExisted = newCustomError("存在无效的投票对象")
	ErrInvalidTicket          = newCustomError("Ticket 无效")
	ErrTicketRequired         = newCustomError("需要 Ticket")
	ErrTicketUsageExceed      = newCustomError("Ticket 使用超过上限")
	ErrCaptchaRequired        = newCustomError("需人机验证")
	ErrCaptchaInvalid         = newCustomError("人机验证失败")
	ErrServerInternal         = newCustomError("服务器内部错误")
	ErrDataLoad               = newCustomError("数据读取异常")
	ErrDataUpdate             = newCustomError("数据更新异常")
	ErrGetTicket              = newCustomError("获取 Ticket 失败")
)
