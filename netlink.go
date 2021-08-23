package netlink

import (
	"golang.org/x/sys/unix"
)

const byteOrder = unix.LittleEndian

type NetlinkMessage struct {
	Header unix.NlMsghdr
	Data []byte
}

func serializeNetlinkMessage(m NetlinkMessage) ([]byte, error) {
	buffer := make([]byte, m.Len)
	byteOrder.PutUint32(buffer[:4], m.Header.Len)
	byteOrder.PutUint16(buffer[4:6], m.Header.Type)
	byteOrder.PutUint16(buffer[6:8], m.Header.Flags)
	byteOrder.PutUint32(buffer[8:12], m.Header.Seq)
	byteOrder.PutUint32(buffer[12:16], m.Header.Pid)
	Copy(buffer[16:], m.Header.Data)
	return buffer
}