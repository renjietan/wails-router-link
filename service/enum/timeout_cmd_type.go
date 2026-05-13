package enum

type TimeoutCmdType = string

const (
	Set     TimeoutCmdType = "set"
	Reset   TimeoutCmdType = "reset"
	Stop    TimeoutCmdType = "stop"
	StopAll TimeoutCmdType = "stop_all"
	Expired TimeoutCmdType = "expired"
	Exists  TimeoutCmdType = "exists"
	Get     TimeoutCmdType = "get"
	List    TimeoutCmdType = "list"
	ListAll TimeoutCmdType = "list_all"
	Stats   TimeoutCmdType = "stats"
	Quit    TimeoutCmdType = "quit"
)
