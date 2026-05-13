package udp_utils_gen

import (
	"github.com/elliotchance/orderedmap/v3"
	"wails-router-link/service/api/udp/utils/parse"
	"wails-router-link/service/api/udp/utils/stuct"
)

func Login() (res []byte) {
	m := orderedmap.NewOrderedMap[string, udp_utils_struct.Base]()
	m.Set("buf", udp_utils_struct.Base{
		Data:  "test_login",
		Range: []int{},
		Size:  0,
	})
	res = udp_utils_parse.ParseSendBuf(m)
	return res
}
