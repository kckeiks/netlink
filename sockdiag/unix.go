package sockdiag

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/kckeiks/netlink"
	"golang.org/x/sys/unix"
)

const (
	UnixDiagReqLen      = 24
	NlUnixDiagReqMsgLen = 40 // includes netlink header
	NlUnixDiagMsgLen    = 40
)

var UnixMsgLenError = errors.New("unix: invalid msg length")

type UnixDiagReq struct {
	Family   uint8
	Protocol uint8
	Pad      uint16
	States   uint32
	Inode    uint32
	Show     uint32
	Cookie   [2]uint32
}

type UnixDiagMsg struct {
	Family uint8
	Type   uint8
	State  uint8
	Pad    uint8
	Inode  uint32
	Cookie [2]uint32
}

func NewUnixNetlinkMsg(nlHeader unix.NlMsghdr, unixHeader UnixDiagReq) ([]byte, error) {
	if nlHeader.Len != NlUnixDiagReqMsgLen {
		return nil, UnixMsgLenError
	}
	msg := netlink.NewSerializedNetlinkMessage(nlHeader)
	ih, err := SerializeUnixDiagReq(unixHeader)
	if err != nil {
		return nil, err
	}
	copy(msg[unix.NLMSG_HDRLEN:], ih)
	return msg, nil
}

func SerializeUnixDiagReq(req UnixDiagReq) ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, UnixDiagReqLen))
	b.Reset()
	err := binary.Write(b, netlink.ByteOrder, req)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DeserializeUnixDiagMsg(data []byte) (*UnixDiagMsg, error) {
	msg := UnixDiagMsg{}
	b := bytes.NewBuffer(data)
	err := binary.Read(b, netlink.ByteOrder, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
