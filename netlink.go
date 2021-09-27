package netlink

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/unix"
	"os"
	"errors"
)

var NlmsgAlignTo uint32 = 4
var OSPageSize = os.Getpagesize()
var NlMsgDoesNotFit = errors.New("nlmsg does not fit into buffer")
var NlMsgHeaderError = errors.New("nlmsghdr error")



type NetlinkMessage struct {
	Header unix.NlMsghdr
	Payload []byte
}

// Round the length of a netlink message
func nlmAlignOf(msglen uint32) uint32 {
	return (msglen + NlmsgAlignTo - 1) & ^(NlmsgAlignTo - 1)
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

func DeserializeNetlinkMsg(data []byte) (*NetlinkMessage, error) {
	len := nlmAlignOf(ByteOrder.Uint32(data[:4]))
	if !IsOkToDeserialize(data, len) {
		return nil, NlMsgDoesNotFit
	}
	serializedData := bytes.NewBuffer(data[:unix.NLMSG_HDRLEN])
	h := unix.NlMsghdr{}
	err := binary.Read(serializedData, ByteOrder, &h)
	if err != nil {
		return nil, err
	}
	return &NetlinkMessage{Header: h, Payload: data[unix.NLMSG_HDRLEN:len]}, nil
}

func IsOkToDeserialize(data []byte, nlmsglen uint32) bool {
	bufLen := uint32(len(data))
	return unix.NLMSG_HDRLEN <= bufLen && nlmsglen >= unix.NLMSG_HDRLEN && bufLen >= nlmsglen
}

func ParseNetlinkMessage(data []byte) ([]NetlinkMessage, error) {
	nlmsgs := make([]NetlinkMessage, 0)
	for len(data) >= unix.NLMSG_HDRLEN {
		len := ByteOrder.Uint32(data[:4])
		msg, err := DeserializeNetlinkMsg(data[:len])
		if err != nil {
			return nil, err
		}
		nlmsgs = append(nlmsgs, *msg)
		data = data[len:]
	}
	return nlmsgs, nil
}

func ReceiveMessage(fd int) (*NetlinkMessage, error) {
	b := make([]byte, OSPageSize)
	n, _, _ := unix.Recvfrom(fd, b, 0)
	return DeserializeNetlinkMsg(b[:n])
}

func ReceiveNetlinkMessage(fd int) ([]NetlinkMessage, error) {
	nlmsgs := make([]NetlinkMessage, 0)
	for done := false; !done; {
		b := make([]byte, OSPageSize)
		r, _, _ := unix.Recvfrom(fd, b, 0)
		if r == 0 {
			return nlmsgs, nil
		}
		parsedMsgs, err := ParseNetlinkMessage(b[:r])
		if err != nil {
			return nil, err
		}
		for _, msg := range parsedMsgs {
			if msg.Header.Type == unix.NLMSG_DONE {
				done = true
				break
			}
			if msg.Header.Type == unix.NLMSG_ERROR {
				return nil, NlMsgHeaderError
			}
			nlmsgs = append(nlmsgs, msg)
		}
	}
	return nlmsgs, nil
}
