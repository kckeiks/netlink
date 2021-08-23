package netlink

import (
	"golang.org/x/sys/unix"
)


const SOCK_DIAG_BY_FAMILY = 20

type NetlinkMessage struct {
	Header unix.NlMsghdr,
	Data []byte
}

type InetDiagSockID struct {
	SPort  [2]byte   // source port          __be16  idiag_sport;
	DPort  [2]byte   // destination port     __be16  idiag_dport;
	Src    [16]byte  // source address       __be32  idiag_src[4];
	Dst    [16]byte  // destination address  __be32  idiag_dst[4];
	If     uint32
	Cookie [2]uint32
}

// inet request structure
type InetDiagReqV2 struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	Pad      uint8
	States   uint32
	ID       InetDiagSockID
}

// inet response structure
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
