package dto

import "wails-router-link/service/enum"

type WS_META_File struct {
	EVENT enum.WS_EVENT `json:"event"`
	TYPE  enum.WS_TYPE  `json:"type"`
	DATA  any           `json:"data"`
}
