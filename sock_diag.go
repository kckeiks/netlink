package netlink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/unix"
)

const sizeOfMessageWithInetDiagReqV2 = 72
const sizeOfInetDiagReqV2 = 56
const SOCK_DIAG_BY_FAMILY = 20

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

func serializeInetDiagReqV2(req InetDiagReqV2) []byte {
	b := bytes.NewBuffer(make([]byte, sizeOfInetDiagReqV2))
	b.Reset()
	err := binary.Write(b, byteOrder, req)
	if err != nil {
		panic("Error: failed to serialize InetDiagReqV2.")
	}
	return b.Bytes()
}

func GetInetDiagMsg() error  {
	// have:
	// m.Header
	// m.Data
	// algoright:
	//  create socket conn (or FD??)
	// serialize data to be sent
	// using fd, use SendTo to send query (figure arguments)
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_SOCK_DIAG)
	if err != nil {
		fmt.Println("Error creating socket.")
		return err
	}

	inetReq := InetDiagReqV2{
		Family: unix.AF_INET,
		Protocol: unix.IPPROTO_TCP,
		States: ^uint32(0),
	}

	header := unix.NlMsghdr{
		Len: sizeOfMessageWithInetDiagReqV2,
		Type: SOCK_DIAG_BY_FAMILY,
		Flags: (unix.NLM_F_REQUEST | unix.NLM_F_DUMP),
		Pid: 0,
	}

	m := NetlinkMessage{
		Header: header,
		Data: serializeInetDiagReqV2(inetReq),
	}

	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	unix.Sendto(fd, SerializeNetlinkMessage(m), 0, addr)

	fmt.Printf("%+v\n", m)
	return nil
}