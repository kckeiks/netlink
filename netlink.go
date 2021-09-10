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
	Payload []byte
}

func NewSerializedNetlinkMessage(h unix.NlMsghdr) []byte {
	b := make([]byte, h.Len)
	ByteOrder.PutUint32(b[:4], h.Len)
	ByteOrder.PutUint16(b[4:6], h.Type)
	ByteOrder.PutUint16(b[6:8], h.Flags)
	ByteOrder.PutUint32(b[8:12], h.Seq)
	ByteOrder.PutUint32(b[12:16], h.Pid)
	return b
}

func DeserializeNetlinkMsg(data []byte) NetlinkMessage {
	if len(data) < unix.SizeofNlMsghdr {
		panic("Error: Could not deserialize. Invalid length for serialized NlMsghdr.")
	}
	serializedData := bytes.NewBuffer(data[:unix.SizeofNlMsghdr])
	h := unix.NlMsghdr{}

	err := binary.Read(serializedData, ByteOrder, &h)
	if err != nil {
		panic("Error: Could not deserialize NlMsghdr.")
	}

	return NetlinkMessage{Header: h, Payload: data[unix.SizeofNlMsghdr:]}
}

func ParseNetlinkMessage(data []byte) []NetlinkMessage {
	var msgs []NetlinkMessage
	for len(data) > unix.NLMSG_HDRLEN {
		l := ByteOrder.Uint32(data[:4])
		nlmsg := DeserializeNetlinkMsg(data[:l])
		msgs = append(msgs, nlmsg)
		data = data[l:]
	}
	return msgs
}

func ReceiveMessage(fd int) NetlinkMessage {
	b := make([]byte, OSPageSize)
	n, _, _ := unix.Recvfrom(fd, b, 0)
	return DeserializeNetlinkMsg(b[:n]) 
}

func ReceiveMultipartMessage(fd int) []NetlinkMessage{
	var msgs []NetlinkMessage
	for done := false; !done; {
		b := make([]byte, OSPageSize)
		n, _, _ := unix.Recvfrom(fd, b, 0)
		for _, msg := range ParseNetlinkMessage(b[:n]) {
			if msg.Header.Type == unix.NLMSG_DONE {
				done = true
				break
			}
			msgs = append(msgs, msg) 
		}
		
	}
	return msgs 
}