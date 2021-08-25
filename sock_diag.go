package netlink

import (
	"fmt"
	"golang.org/x/sys/unix"
)

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
	b := make([]byte, sizeOfInetDiagReqV2)
	b[0] = req.Family
	b[1] = req.Protocol
	b[2] = req.Ext
	b[3] = req.Pad
	byteOrder.PutUint32(b[4:8], req.States)
	copy(b[8:10], req.ID.SPort[:])
	copy(b[10:12], req.ID.DPort[:])
	copy(b[12:28], req.ID.Src[:])
	copy(b[28:44], req.ID.Dst[:])
	byteOrder.PutUint32(b[44:48], req.ID.If)
	byteOrder.PutUint32(b[48:52], req.ID.Cookie[0])
	byteOrder.PutUint32(b[52:56], req.ID.Cookie[1])

	return b
}

func SendQuery(m NetlinkMessage) error {
	// have:
	// m.Header
	// m.Data
	// algoright:
	//  create socket conn (or FD??)
	// serialize data to be sent
	// using fd, use SendTo to send query (figure arguments)
	_, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_SOCK_DIAG)
	if err != nil {
		fmt.Println("Error creating socket.")
		return err
	}
	fmt.Println("SUCCESS!")
	return nil
}