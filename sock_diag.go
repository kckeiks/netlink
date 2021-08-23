package netlink

import (
	"fmt"
	"golang.org/x/sys/unix"
)


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