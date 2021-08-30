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

func SerializeInetDiagReqV2(req InetDiagReqV2) []byte {
	b := bytes.NewBuffer(make([]byte, sizeOfInetDiagReqV2))
	b.Reset()
	err := binary.Write(b, byteOrder, req)
	if err != nil {
		panic("Error: failed to serialize InetDiagReqV2.")
	}
	return b.Bytes()
}

func DeserializeInetDiagReqV2(data []byte) InetDiagReqV2 {
	b := bytes.NewBuffer(data)
	req := InetDiagReqV2{}
	err := binary.Read(b, byteOrder, &req)
	if err != nil {
		panic("Error: Could not deserialize InetDiagReqV2.")
	}
	return req
}

func GetInetDiagMsg() error  {
	// Algorithm:
	//  create socket
	//  init data
	// 	serialize data to be sent
	//  use SendTo to send query
	//  use Rcvdfrom to get response

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

	h := unix.NlMsghdr{
		Len: sizeOfMessageWithInetDiagReqV2,
		Type: SOCK_DIAG_BY_FAMILY,
		Flags: (unix.NLM_F_REQUEST | unix.NLM_F_DUMP),
		Pid: 0,
	}

	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	unix.Sendto(fd, NewSerializedNetlinkMsg(h, SerializeInetDiagReqV2(inetReq)), 0, addr)

	readBuffer := make([]byte, OSPageSize)
	n, _, _ := unix.Recvfrom(fd, readBuffer, 0)

	readBuffer = readBuffer[:n]
	for _, msg := range ParseNetlinkMessages(readBuffer) {
		fmt.Printf("Header: %+v\n", msg.Header)
		fmt.Printf("Value: %+v\n", msg.Data)
		fmt.Println("-------")
	}
	return nil
}