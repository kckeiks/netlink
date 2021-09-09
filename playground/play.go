package main

import (
    "fmt"
    "github.com/kckeiks/netlink"
    "golang.org/x/sys/unix"
)

func main() {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_SOCK_DIAG)
	if err != nil {
		panic("Error creating socket.")
	}

	inetReq := netlink.InetDiagReqV2{
		Family: unix.AF_INET,
		Protocol: unix.IPPROTO_TCP,
		States: ^uint32(0),
	}

	h := unix.NlMsghdr{
		Len: netlink.SizeOfMessageWithInetDiagReqV2,
		Type: netlink.SOCK_DIAG_BY_FAMILY,
		Flags: (unix.NLM_F_REQUEST | unix.NLM_F_DUMP),
		Pid: 0,
	}

	inetReq.CreateInetNetlinkMsg(h)


	msg := NewDecodedInetNetlinkMsg(h, inetheader)
	msg := NewDecodedLinuxNetlinkMsg(h, inetheader)

	response := SendNetlinkMessage(msg)

	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	unix.Sendto(fd, netlink.NewEncodedNetlinkMsg(h, netlink.SerializeInetDiagReqV2(inetReq)), 0, addr)

	readBuffer := make([]byte, netlink.OSPageSize)
	n, _, _ := unix.Recvfrom(fd, readBuffer, 0)

	readBuffer = readBuffer[:n]
	for _, msg := range netlink.ParseNetlinkMsgs(readBuffer) {
		fmt.Printf("Header: %+v\n", msg.Header)
		fmt.Printf("Value: %+v\n", netlink.DeserializeInetDiagMsg(msg.Data))
		fmt.Println("-------")
	}
}