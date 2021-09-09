package main

import (
    "fmt"
    // "github.com/kckeiks/netlink"
	"github.com/kckeiks/netlink/sockdiag"
    "golang.org/x/sys/unix"
)

func main() {
	inetReq := sockdiag.InetDiagReqV2{
		Family: unix.AF_INET,
		Protocol: unix.IPPROTO_TCP,
		States: ^uint32(0),
	}

	h := unix.NlMsghdr{
		Len: sockdiag.SizeOfMessageWithInetDiagReqV2,
		Type: sockdiag.SOCK_DIAG_BY_FAMILY,
		Flags: (unix.NLM_F_REQUEST | unix.NLM_F_DUMP),
		Pid: 0,
	}

	nlmsg := sockdiag.NewInetNetlinkMsg(h, inetReq)

	response := sockdiag.SendInetQuery(nlmsg)

	for _, msg := range response {
		fmt.Printf("InetDiagMsg: %+v\n", msg)
	}
}