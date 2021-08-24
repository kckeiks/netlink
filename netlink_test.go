package netlink

import (
	"testing"
	"golang.org/x/sys/unix"
	""
)

func TestValuesWhenSerializeNetlinkMessage(t *testing.T) {
	// Given: a nl message with some values
	m := NetlinkMessage{}
	data := [4]byte{0, 0, 0, 0}
	m.Data = data[:]
	header := unix.NlMsghdr{}
	header.Len = uint32(unix.SizeofNlMsghdr + len(data))
	header.Type = uint16(2)
	header.Flags = uint16(5)
	header.Seq = uint32(6)
	header.Pid = uint32(11)
	m.Header = header
	// When: we serialize the msg
	result := SerializeNetlinkMessage(m)
	// Then: we get the correct number of bytes
	

}

func TestLengthWhenSerializeNetlinkMessage(t *testing.T) {
	// Given: a nl message
	m := NetlinkMessage{}
	header := unix.NlMsghdr{}
	data := [4]byte{0, 0, 0, 0}
	header.Len = uint32(unix.SizeofNlMsghdr + len(data))
	m.Header = header
	m.Data = data[:]
	// When: we serialize the msg
	result := SerializeNetlinkMessage(m)
	// Then: we get the correct number of bytes
	if len(result) != 20 {
		t.Fatalf("Failed. Incorrect length.")
	}
}