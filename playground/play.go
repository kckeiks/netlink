package main

import (
    "fmt"
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
		Len: sockdiag.NL_INET_DIAG_REQ_V2_MSG_LEN,
		Type: sockdiag.SOCK_DIAG_BY_FAMILY,
		Flags: (unix.NLM_F_REQUEST | unix.NLM_F_DUMP),
		Pid: 0,
	}

	nlmsg := sockdiag.NewInetNetlinkMsg(h, inetReq)

	result, _ := sockdiag.SendInetMessage(nlmsg)

	for _, msg := range result {
		fmt.Printf("InetDiagMsg: %+v\n", msg)
	}
}