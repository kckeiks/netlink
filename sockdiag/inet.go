package sockdiag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"golang.org/x/sys/unix"
	"github.com/kckeiks/netlink"
)

const (
	InetDiagReqV2Len        = 56
	NlInetDiagReqV2MsgLen   = 72 // includes netlink header
	NlInetDiagMsgLen        = 72
)

var InetMsgLenError = errors.New("inet: invalid msg length")

type InetDiagSockID struct {
	SPort  [2]byte    // source port          __be16  idiag_sport;
	DPort  [2]byte    // destination port     __be16  idiag_dport;
	Src    [16]byte   // source address       __be32  idiag_src[4];
	Dst    [16]byte   // destination address  __be32  idiag_dst[4];
	If     uint32
	Cookie [2]uint32
}

type InetDiagReqV2 struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	Pad      uint8
	States   uint32
	ID       InetDiagSockID
}

type InetDiagMsg struct {
	Family  uint8
	State   uint8
	Timer   uint8
	Retrans uint8
	ID InetDiagSockID
	Expires uint32
	RQueue  uint32
	WQueue  uint32
	UID     uint32
	Inode   uint32
}

func NewInetNetlinkMsg(nlHeader unix.NlMsghdr, inetHeader InetDiagReqV2) ([]byte, error) {
	if nlHeader.Len != NlInetDiagReqV2MsgLen {
		return nil, InetMsgLenError
	}
	msg := netlink.NewSerializedNetlinkMessage(nlHeader)
	ih, err := SerializeInetDiagReqV2(inetHeader)
	if err != nil {
		return nil, err
	}
	copy(msg[unix.NLMSG_HDRLEN :], ih)
	return msg, nil
}

func SerializeInetDiagReqV2(req InetDiagReqV2) ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, InetDiagReqV2Len))
	b.Reset()
	err := binary.Write(b, netlink.ByteOrder, req)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DeserializeInetDiagMsg(data []byte) (*InetDiagMsg, error) {
	msg := InetDiagMsg{}
	b := bytes.NewBuffer(data)
	err := binary.Read(b, netlink.ByteOrder, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil 
}
