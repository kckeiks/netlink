package netlink

import (
	"golang.org/x/sys/unix"
)

type NetlinkMessage struct {
	Header unix.NlMsghdr
	Data []byte
}
