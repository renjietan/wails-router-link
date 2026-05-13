package udp_utils_parse

import (
	"fmt"
	"slices"

	"github.com/elliotchance/orderedmap/v3"
	"wails-router-link/service/api/udp/utils/stuct"

	"wails-router-link/service/utility"
)

func ParseSendBuf(data *orderedmap.OrderedMap[string, udp_utils_struct.Base]) (res []byte) {
	pData := make(map[string][]byte)

	sumKey := ""
	sumStart := false
	sumTotal := 0
	sumRange := []int{}
	order_keys := slices.Collect(data.Keys())
	for index, key := range order_keys {
		item, _ := data.Get(key)
		if len(item.Range) > 1 {
			sumKey = key
			sumStart = true
			sumTotal = 0
			sumRange = []int{}
			sumRange = append(sumRange, item.Range...)
			pData[key] = make([]byte, item.Size)
			// res = append(res, pData[key]...)
			continue
		}

		var dataBytes []byte

		switch v := item.Data.(type) {
		case int:
			dataBytes = utility.ByteToMuiByte(v, item.Size, "big")
		case []byte:
			if item.Size > 0 {
				if len(v) >= item.Size {
					dataBytes = v[:item.Size]
				} else {
					dataBytes = make([]byte, item.Size)
					copy(dataBytes, v)
				}
			} else {
				dataBytes = v
			}
		case string:
			dataBytes = []byte(v)
			item.Size = len(dataBytes)
		default:
			dataBytes = []byte{}
		}
		pData[key] = dataBytes
		// res = append(res, pData[key]...)
		if sumStart && sumKey != "" && index >= sumRange[0] && index <= sumRange[1] {
			if sumKey == "length2" {
				fmt.Printf("sum: %d, len(dataBytes): %d\n", sumTotal, len(dataBytes))
			}
			sumTotal += len(dataBytes)
			if index == sumRange[1] {
				item, _ := data.Get(sumKey)
				pData[sumKey] = utility.ByteToMuiByte(sumTotal, item.Size, "big")
			}
		}
	}
	for _, key := range order_keys {
		res = append(res, pData[key]...)
		fmt.Printf("%s: %v %#v(length: %d)\n", key, pData[key], pData[key], len(pData[key]))
	}
	return res
}

func ParseRevBuf(data *orderedmap.OrderedMap[string, udp_utils_struct.Base], buf []byte) {
	keys := slices.Collect(data.Keys())
	count := 0
	for _, v := range keys {
		item, _ := data.Get(v)
		item.Data = utility.Tern(
			item.Size > 0,
			buf[count:count+item.Size], buf[count:],
		)
		data.Set(v, item)
		count = count + item.Size
	}
	for k, v := range data.AllFromFront() {
		// fmt.Printf("%s===============%#v==================%v\n", k, v, v)
		fmt.Printf("%s\n", k)
		fmt.Printf("原始值: %#v\n", v.Data)
		fmt.Printf("字符串值: %s\n", string(v.Data.([]byte)))
	}
}
