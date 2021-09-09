package netlink

import (
	"bytes"
	"reflect"
	"testing"
	"golang.org/x/sys/unix"
)

func NewTestNlMsghdr() unix.NlMsghdr {
	h := unix.NlMsghdr{}
	h.Len = unix.SizeofNlMsghdr
	h.Type = 2
	h.Flags = 5
	h.Seq = 6
	h.Pid = 11
	return h
}

func newTestSerializedNetlinkMsg(h unix.NlMsghdr, data []byte) []byte {
	if h.Len != (uint32(len(data)) + unix.SizeofNlMsghdr) {
		panic("Error: Invalid NlMsghdr.Len.")
	}
	b := make([]byte, h.Len)
	byteOrder.PutUint32(b[:4], h.Len)
	byteOrder.PutUint16(b[4:6], h.Type)
	byteOrder.PutUint16(b[6:8], h.Flags)
	byteOrder.PutUint32(b[8:12], h.Seq)
	byteOrder.PutUint32(b[12:16], h.Pid)
	copy(b[16:], data)
	return b
}

func TestNewNetlinkMessage(t *testing.T) {
	// Given: a NlMsghdr header
	h := NewTestNlMsghdr()
	// Given: length of nl msg with 4 more bytes of space 
	h.Len = 16 + 4

	// When: we serialize the header and the data
	serializedData := NewNetlinkMessage(h)

	// Then: the message was serialized with the correct data
	if h.Len != testByteOrder.Uint32(serializedData[:4]) {
		t.Fatalf("NlMsghdr.Length = %d, expected %d", testByteOrder.Uint32(serializedData[:4]), h.Len)
	}
	if h.Type != testByteOrder.Uint16(serializedData[4:6]) {
		t.Fatalf("NlMsghdr.Type = %d, expected %d", testByteOrder.Uint16(serializedData[4:6]), h.Type)
	}
	if h.Flags != testByteOrder.Uint16(serializedData[6:8]) {
		t.Fatalf("NlMsghdr.Flags = %d, expected %d", testByteOrder.Uint16(serializedData[6:8]), h.Flags)
	}
	if h.Seq != testByteOrder.Uint32(serializedData[8:12]) {
		t.Fatalf("NlMsghdr.Seq = %d, expected %d", testByteOrder.Uint32(serializedData[8:12]), h.Seq)
	}
	if h.Pid != testByteOrder.Uint32(serializedData[12:16]) {
		t.Fatalf("NlMsghdr.Pid = %d, expected %d", testByteOrder.Uint32(serializedData[12:16]), h.Pid)
	}
	if len(serializedData[16:]) != 4 {
		t.Fatalf("Len(serializedData) = %d, expected %d", len(serializedData[16:]), 4)
	}
}

func TestDeserializeNetlinkMsg(t *testing.T) {
	// Given: a serialized netlink message
	h := NewTestNlMsghdr()
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	h.Len = h.Len + uint32(len(data))
	serializedData := newTestSerializedNetlinkMsg(h, data[:])
	
	// When: we deserialize the message
	nlmsg := DeserializeNetlinkMsg(serializedData)

	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(nlmsg.Header, h) {
		t.Fatalf("Given NlMsghdr %+v and deserialized is %+v,", nlmsg.Header, h)
	}
	// Then: the extra data was returned
	if bytes.Compare(nlmsg.Data, data[:]) != 0 {
		t.Fatalf("Extra data=%d, expected %d", nlmsg.Data, data)
	}
}

func TestDeserializeNetlinkMsgWithOutData(t *testing.T) {
	// Given: a serialized netlink message without extra data
	h := NewTestNlMsghdr()
	data := []byte{} 
	h.Len = uint32(unix.SizeofNlMsghdr)
	serializedData := newTestSerializedNetlinkMsg(h, data)
	
	// When: we deserialize the message
	nlmsg := DeserializeNetlinkMsg(serializedData)

	// Then: empty slice is returned for payload
	if len(nlmsg.Data) != 0 {
		t.Fatalf("Extra data=%d, expected [].", nlmsg.Data)
	}
}

func TestDeserializeNetlinkMsgBadLen(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Error: did not panic.")
        }
    }()
	// Given: a serialized nl message with extra data but 
	// we do not update length in header
	h := NewTestNlMsghdr()
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	serializedData := newTestSerializedNetlinkMsg(h, data[:])
	// When: we deserialize the message
	// Then: we panic
	DeserializeNetlinkMsg(serializedData)
}

func TestParseNetlinkMessages(t *testing.T) {
	// Given: a list of serialized netlink messages
	h1 := NewTestNlMsghdr()
	data1 := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	h1.Len = h1.Len + uint32(len(data1))
	h2 := NewTestNlMsghdr()
	data2 := [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x01}
	h2.Len = h2.Len + uint32(len(data2))

	var nlmsgs []byte
	nlmsgs = append(nlmsgs, newTestSerializedNetlinkMsg(h1, data1[:])...)
	nlmsgs = append(nlmsgs, newTestSerializedNetlinkMsg(h2, data2[:])...)

	// When: parse these serialized data
	result := ParseNetlinkMessages(nlmsgs)

	// Then: We get the messages as expected
	expectedNlMsg1 := NetlinkMessage{Header: h1, Data: data1[:]}
	expectedNlMsg2 := NetlinkMessage{Header: h2, Data: data2[:]}

	if !reflect.DeepEqual(result[0], expectedNlMsg1) {
		t.Fatalf("Given first Netlink Msg %+v but received %+v,", result[0], expectedNlMsg1)
	}
	if !reflect.DeepEqual(result[1], expectedNlMsg2) {
		t.Fatalf("Given second Netlink Msg %+v but received %+v,", result[1], expectedNlMsg2)
	}
}