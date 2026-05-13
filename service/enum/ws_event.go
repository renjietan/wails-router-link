package enum

type WS_EVENT = string

const (
	WS_EVENT_UNKNOWN WS_EVENT = "unknown"
	WS_EVENT_PING    WS_EVENT = "ping"
	WS_EVENT_PONG    WS_EVENT = "pong"
	WS_EVENT_LOGIN   WS_EVENT = "login"
	WS_EVENT_FILE    WS_EVENT = "file"
)
