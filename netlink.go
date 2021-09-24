package netlink

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/sys/unix"
	"os"
)

var OSPageSize = os.Getpagesize()

type NetlinkMessage struct {
	Header unix.NlMsghdr
	Payload []byte
}

// Round the length of a netlink message
func nlmAlignOf(msglen int) int {
	return (msglen + unix.NLMSG_ALIGNTO - 1) & ^(unix.NLMSG_ALIGNTO - 1)
}

func NewSerializedNetlinkMessage(h unix.NlMsghdr) []byte {
	b := make([]byte, h.Len)
	ByteOrder.PutUint32(b[:4], h.Len)
	ByteOrder.PutUint16(b[4:6], h.Type)
	ByteOrder.PutUint16(b[6:8], h.Flags)
	ByteOrder.PutUint32(b[8:12], h.Seq)
	ByteOrder.PutUint32(b[12:16], h.Pid)
	return b
}

func DeserializeNetlinkMsg(data []byte) NetlinkMessage {
	l := nlmAlignOf(int(ByteOrder.Uint32(data[:4])))
	if len(data) < unix.NLMSG_HDRLEN || l > len(data) {
		panic("Error: Could not deserialize. Invalid length for serialized NlMsghdr.")
	}
	serializedData := bytes.NewBuffer(data[:unix.NLMSG_HDRLEN ])
	h := unix.NlMsghdr{}

	err := binary.Read(serializedData, ByteOrder, &h)
	if err != nil {
		panic("Error: Could not deserialize NlMsghdr.")
	}

	return NetlinkMessage{Header: h, Payload: data[unix.NLMSG_HDRLEN :]}
}

func ParseNetlinkMessage(data []byte) []NetlinkMessage {
	var msgs []NetlinkMessage
	for len(data) > unix.NLMSG_HDRLEN {
		l := ByteOrder.Uint32(data[:4])
		nlmsg := DeserializeNetlinkMsg(data[:l])
		msgs = append(msgs, nlmsg)
		data = data[l:]
	}
	return msgs
}

func ReceiveMessage(fd int) NetlinkMessage {
	b := make([]byte, OSPageSize)
	n, _, _ := unix.Recvfrom(fd, b, 0)
	return DeserializeNetlinkMsg(b[:n]) 
}

func ReceiveMultipartMessage(fd int) []NetlinkMessage{
	var msgs []NetlinkMessage
	for done := false; !done; {
		b := make([]byte, OSPageSize)
		n, _, _ := unix.Recvfrom(fd, b, 0)
		// TODO: Check that this does not return 0
		for {
			// msg, err :=  getNextNetlingMessage where err could be:
			// len(b) <= nlmsghdrlen || or desrializer isse
			// basically issues that can happen when desiarilising 
			//	l := ByteOrder.Uint32(b[:4])
			//  return DeserializeNetlinkMsg(b[:l])

			// e.g. ParseNetlinkMessage doesnt add value
			
			if msg.Header.Type == unix.NLMSG_DONE {
				done = true
				break
			}
			msgs = append(msgs, msg) 
		}
		
	}
	return msgs 
}

func ReceiveNetlinkMessage(fd int) []NetlinkMessage{
	nlmsgs := make([]NetlinkMessage, 0)
	buf := make([]byte, OSPageSize)
	r, _, _ := unix.Recvfrom(fd, buf, 0)
	buf = buf[:r]
	if len(buf) <= unix.NLMSG_HDRLEN {
		panic("Error: Invalid length of first Nl MSG.")
	}
	firstMsgLen := ByteOrder.Uint32(buf[:4])
	firstMsg := DeserializeNetlinkMsg(buf[:firstMsgLen])
	buf = buf[firstMsgLen:]
	nlmsgs = append(nlmsgs, firstMsg)
	if firstMsg.Header.Flags != unix.NLM_F_MULTI {
		return nlmsgs
	}
	// Handle multi-part message
	responseMsgs := ParseNetlinkMessage(buf)
	nlmsgs = append(nlmsgs, responseMsgs...)
	responseMsgs = ReceiveMultipartMessage(fd)
	return append(nlmsgs, responseMsgs...)
}


// TODO: this should do some validation on length of b
func getNextNetlingMessage(b []byte) (NetlinkMessage, uint32, error) {
	l := ByteOrder.Uint32(b[:4])
	return DeserializeNetlinkMsg(b[:l]), l, nil
}
