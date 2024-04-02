package errors

type customErr struct {
	Msg string
}

func newError(msg string) customErr {
	return customErr{
		Msg: msg,
	}
}

func (e *customErr) WithDetail(msg string) *customErr {
	return &customErr{
		Msg: e.Msg + ": " + msg,
	}
}

func (e customErr) Error() string {
	return e.Msg
}

var (
	InvalidParams          = newError("参数错误")
	TicketNeeded           = newError("需要 Ticket")
	TicketInvalid          = newError("Ticket 无效")
	TicketUsageLimitExceed = newError("Ticket 使用超过上限")
	UserNotFound           = newError("用户不存在")
	ServerInternal         = newError("服务器内部错误")
	GetTicketFailed        = newError("获取 Ticket 失败")
)
