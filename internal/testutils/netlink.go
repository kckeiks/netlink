package testutils

import (
	"encoding/binary"

	"golang.org/x/sys/unix"
)

var TestByteOrder = binary.LittleEndian

func NewTestNlMsghdr() unix.NlMsghdr {
	h := unix.NlMsghdr{}
	h.Len = unix.NLMSG_HDRLEN
	h.Type = 0
	h.Flags = 5
	h.Seq = 6
	h.Pid = 11
	return h
}

func NewTestSerializedNetlinkMsg(h unix.NlMsghdr, data []byte) []byte {
	if h.Len != (uint32(len(data)) + unix.NLMSG_HDRLEN) {
		panic("Error: Invalid NlMsghdr.Len.")
	}
	b := make([]byte, h.Len)
	TestByteOrder.PutUint32(b[:4], h.Len)
	TestByteOrder.PutUint16(b[4:6], h.Type)
	TestByteOrder.PutUint16(b[6:8], h.Flags)
	TestByteOrder.PutUint32(b[8:12], h.Seq)
	TestByteOrder.PutUint32(b[12:16], h.Pid)
	copy(b[16:], data)
	return b
}
