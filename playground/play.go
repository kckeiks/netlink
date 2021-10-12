package main

import (
	"fmt"
	"os"

	"github.com/kckeiks/netlink"
	"github.com/kckeiks/netlink/sockdiag"
	"golang.org/x/sys/unix"
)

func SendInetMessage(nlmsg []byte) ([]sockdiag.InetDiagMsg, error) {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_SOCK_DIAG)
	if err != nil {
		panic("Error creating socket.")
	}
	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	unix.Sendto(fd, nlmsg, 0, addr)
	nlmsgs, err := ReceiveNetlinkMessage(fd)
	if err != nil {
		return nil, err
	}
	var idmsgs []sockdiag.InetDiagMsg
	for _, msg := range nlmsgs {
		inetMsg, _ := sockdiag.DeserializeInetDiagMsg(msg.Payload)
		idmsgs = append(idmsgs, *inetMsg)
	}
	return idmsgs, nil
}

func ReceiveNetlinkMessage(fd int) ([]netlink.NetlinkMessage, error) {
	nlmsgs := make([]netlink.NetlinkMessage, 0)
	for done := false; !done; {
		b := make([]byte, os.Getpagesize())
		r, _, _ := unix.Recvfrom(fd, b, 0)
		if r == 0 {
			return nlmsgs, nil
		}
		parsedMsgs, err := netlink.ParseNetlinkMessage(b[:r])
		if err != nil {
			return nil, err
		}
		for _, msg := range parsedMsgs {
			if msg.Header.Type == unix.NLMSG_DONE {
				done = true
				break
			}
			if msg.Header.Type == unix.NLMSG_ERROR {
				return nil, netlink.NlMsgHeaderError
			}
			nlmsgs = append(nlmsgs, msg)
		}
	}
	return nlmsgs, nil
}

func main() {
	inetReq := sockdiag.InetDiagReqV2{
		Family:   unix.AF_INET,
		Protocol: unix.IPPROTO_TCP,
		States:   ^uint32(0),
	}
	h := unix.NlMsghdr{
		Len:   sockdiag.NlInetDiagReqV2MsgLen,
		Type:  sockdiag.SOCK_DIAG_BY_FAMILY,
		Flags: (unix.NLM_F_REQUEST | unix.NLM_F_DUMP),
		Pid:   0,
	}
	nlmsg, _ := sockdiag.NewInetNetlinkMsg(h, inetReq)
	result, _ := SendInetMessage(nlmsg)
	for _, msg := range result {
		fmt.Printf("InetDiagMsg: %+v\n", msg)
	}
}
