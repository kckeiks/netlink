package netlink

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/unix"
)

type NetlinkMessage struct {
	Header unix.NlMsghdr
	Data []byte
}

func NewSerializedNetlinkMsg(h unix.NlMsghdr, data []byte) []byte {
	if h.Len != (uint32(len(data)) + unix.SizeofNlMsghdr) {
		panic("Error: Invalid NlMsghdr.Len.")
	}
	b := make([]byte, h.Len)
	byteOrder.PutUint32(b[:4], h.Len)
	byteOrder.PutUint16(b[4:6], h.Type)
	byteOrder.PutUint16(b[6:8], h.Flags)
	byteOrder.PutUint32(b[8:12], h.Seq)
	byteOrder.PutUint32(b[12:16], h.Pid)
	copy(b[16:], data)
	return b
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

func deserializeNetlinkMessage(data []byte) (unix.NlMsghdr, []byte) {
	if len(data) < unix.SizeofNlMsghdr {
		panic("Error: Could not deserialize. Invalid length for serialized NlMsghdr.")
	}
	
	serializedData := bytes.NewBuffer(data[:unix.SizeofNlMsghdr])
	header := unix.NlMsghdr{}

	err := binary.Read(serializedData, byteOrder, &header)
	if err != nil {
		panic("Error: Could not deserialize NlMsghdr.")
	}

	if len(data) == unix.SizeofNlMsghdr {
		return header, nil
	}
	
	return header, data[unix.SizeofNlMsghdr:]
}