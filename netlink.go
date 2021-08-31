package netlink

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/unix"
	"os"
)

var OSPageSize = os.Getpagesize()

type NetlinkMessage struct {
	Header unix.NlMsghdr
	Data []byte
}

func NewEncodedNetlinkMsg(h unix.NlMsghdr, data []byte) []byte {
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

func ParseNetlinkMsg(data []byte) (unix.NlMsghdr, []byte) {
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

func ParseNetlinkMsgs(data []byte) []NetlinkMessage {
	var msgs []NetlinkMessage
	for len(data) > unix.NLMSG_HDRLEN {
		l := byteOrder.Uint32(data[:4])
		h, d := ParseNetlinkMsg(data[:l])
		msgs = append(msgs, NetlinkMessage{Header: h, Data: d})
		data = data[l:]
	}
	return msgs
}