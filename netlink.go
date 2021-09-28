package netlink

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/unix"
	"errors"
)

var ByteOrder = binary.LittleEndian

var NlmsgAlignTo uint32 = 4

var (
	NlMsgDoesNotFit = errors.New("nlmsg does not fit into buffer")
	NlMsgHeaderError = errors.New("nlmsghdr error")
)

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
	if !IsOkToParse(data, len) {
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

func IsOkToParse(data []byte, nlmsglen uint32) bool {
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
