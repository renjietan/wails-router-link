package utility

import "encoding/binary"

func ByteToMuiByte(i int, length int, endian string) []byte {
	buf := make([]byte, length)
	switch length {
	case 1:
		buf[0] = byte(i)
	case 2:
		Tern(endian == "big",
			func() { binary.BigEndian.PutUint16(buf, uint16(i)) },
			func() { binary.LittleEndian.PutUint16(buf, uint16(i)) })
	case 4:
		Tern(endian == "big",
			func() { binary.BigEndian.PutUint32(buf, uint32(i)) },
			func() { binary.LittleEndian.PutUint32(buf, uint32(i)) })
	case 8:
		Tern(endian == "big",
			func() { binary.BigEndian.PutUint64(buf, uint64(i)) },
			func() { binary.LittleEndian.PutUint64(buf, uint64(i)) })
	}
	return buf
}
