package utils

import "bytes"

func WriteString(buf *bytes.Buffer, s string) {
	buf.WriteString(s)
	buf.WriteByte(0x00)
}

