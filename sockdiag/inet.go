package sockdiag

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/unix"
	"github.com/kckeiks/netlink"
)

const SizeOfMessageWithInetDiagReqV2 = 72
const sizeOfInetDiagReqV2 = 56
const sizeOfInetDiagMsg = 72

type InetDiagSockID struct {
	SPort  [2]byte   // source port          __be16  idiag_sport;
	DPort  [2]byte    // destination port     __be16  idiag_dport;
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

func NewInetNetlinkMsg(nlh unix.NlMsghdr, inetHeader InetDiagReqV2) []byte {
	if nlh.Len != SizeOfMessageWithInetDiagReqV2 {
		panic("Error: Invalid NlMsghdr.Len.")
	}
	msg := netlink.NewSerializedNetlinkMessage(nlh)
	ih := SerializeInetDiagReqV2(inetHeader)
	copy(msg[unix.SizeofNlMsghdr:], ih)
	return msg
}

func SerializeInetDiagReqV2(req InetDiagReqV2) []byte {
	b := bytes.NewBuffer(make([]byte, sizeOfInetDiagReqV2))
	b.Reset()
	err := binary.Write(b, netlink.ByteOrder, req)
	if err != nil {
		panic("Error: failed to serialize InetDiagReqV2.")
	}
	return b.Bytes()
}

func DeserializeInetDiagMsg(data []byte) InetDiagMsg {
	msg := InetDiagMsg{}
	b := bytes.NewBuffer(data)
	err := binary.Read(b, netlink.ByteOrder, &msg)
	if err != nil {
		panic("Error: Could not parse InetDiagMsg.")
	}
	return msg 
}


func SendInetQuery(nlmsg []byte) []InetDiagMsg {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_SOCK_DIAG)
	if err != nil {
		panic("Error creating socket.")
	}

	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	unix.Sendto(fd, nlmsg, 0, addr)

	nlmsgs := netlink.ReceiveMultipartMessage(fd)

	idmsgs := make([]InetDiagMsg, len(nlmsgs))

	for _, msg := range nlmsgs {
		idmsgs = append(idmsgs, DeserializeInetDiagMsg(msg.Payload))
	}

	return idmsgs
}