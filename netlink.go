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