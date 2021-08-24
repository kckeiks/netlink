package netlink

import (
	"fmt"
	"golang.org/x/sys/unix"
)

const SOCK_DIAG_BY_FAMILY = 20

type InetDiagSockID struct {
	SPort  uint16   // source port          __be16  idiag_sport;
	DPort  uint16   // destination port     __be16  idiag_dport;
	Src    [4]uint32  // source address       __be32  idiag_src[4];
	Dst    [4]uint32  // destination address  __be32  idiag_dst[4];
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
	buffer := make([]byte, 56)
	buffer[0] = req.Family
	buffer[1] = req.Protocol
	buffer[2] = req.Ext
	buffer[3] = req.Pad
	byteOrder.PutUint32(buffer[4:8], req.States)
	ipAddrByteOrder.PutUint16(buffer[8:10], req.ID.SPort)
	ipAddrByteOrder.PutUint16(buffer[10:12], req.ID.DPort)
	ipAddrByteOrder.PutUint32(buffer[12:28], req.ID.Src)
	ipAddrByteOrder.PutUint32(buffer[28:44], req.ID.Dst)
	byteOrder.PutUint32(buffer[44:48], req.ID.If)
	byteOrder.PutUint32(buffer[48:56], req.ID.Cookie)
	return buffer
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