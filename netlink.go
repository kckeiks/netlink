package netlink

import "golang.org/x/sys/unix"

type NetlinkMessage struct {
	Header unix.NlMsghdr
	Data []byte
}

func SerializeNetlinkMessage(m NetlinkMessage) []byte {
	buffer := make([]byte, m.Header.Len)
	byteOrder.PutUint32(buffer[:4], m.Header.Len)
	byteOrder.PutUint16(buffer[4:6], m.Header.Type)
	byteOrder.PutUint16(buffer[6:8], m.Header.Flags)
	byteOrder.PutUint32(buffer[8:12], m.Header.Seq)
	byteOrder.PutUint32(buffer[12:16], m.Header.Pid)
	copy(buffer[16:], m.Data)
	return buffer
}