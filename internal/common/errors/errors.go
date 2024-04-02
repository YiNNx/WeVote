package errors

import "github.com/YiNNx/WeVote/internal/config"

type customErr struct {
	Msg string
}

func newError(msg string) customErr {
	return customErr{
		Msg: msg,
	}
}

func (e *customErr) WithErrDetail(err error) *customErr {
	if !config.C.Server.DebugMode {
		return e
	}
	return &customErr{
		Msg: e.Msg + ": " + e.Error(),
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
	DataLoadFailed         = newError("数据读取异常")
	DataUpdateFailed       = newError("数据更新异常")
	GetTicketFailed        = newError("获取 Ticket 失败")
)
