package netlink

import (
	"bytes"
	"reflect"
	"testing"
	"golang.org/x/sys/unix"
)

func CreateTestNlMsghdr() unix.NlMsghdr {
	h := unix.NlMsghdr{}
	h.Len = unix.SizeofNlMsghdr
	h.Type = 2
	h.Flags = 5
	h.Seq = 6
	h.Pid = 11
	return h
}

func NewEncodedNetlinkMsg(h unix.NlMsghdr, data []byte) []byte {
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

func TestParseNetlinkMsg(t *testing.T) {
	// Given: a serialized netlink message
	h := CreateTestNlMsghdr()
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	h.Len = h.Len + uint32(len(data))
	serializedData := NewEncodedNetlinkMsg(h, data[:])
	
	// When: we deserialize the message
	result, xdata := ParseNetlinkMsg(serializedData)

	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, h) {
		t.Fatalf("Given NlMsghdr %+v and deserialized is %+v,", result, h)
	}
	// Then: the extra data was returned
	if bytes.Compare(xdata, data[:]) != 0 {
		t.Fatalf("Extra data=%d, expected %d", xdata, data)
	}
}

func TestDeserializeNetlinkMessageWithOutData(t *testing.T) {
	// Given: a serialized netlink message without extra data
	h := CreateTestNlMsghdr()
	data := []byte{} 
	h.Len = uint32(unix.SizeofNlMsghdr)
	serializedData := NewEncodedNetlinkMsg(h, data)
	
	// When: we deserialize the message
	_, xdata := ParseNetlinkMsg(serializedData)

	// Then: nil is returned for the extra data
	if xdata != nil {
		t.Fatalf("Extra data=%d, expected nil.", xdata)
	}
}

func TestParseNetlinkMsgBadLen(t *testing.T) {
	defer func() {
        if r := recover(); r == nil {
            t.Errorf("Error: did not panic.")
        }
    }()
	// Given: a serialized nl message with extra data but 
	// we do not update length in header
	h := CreateTestNlMsghdr()
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	serializedData := NewEncodedNetlinkMsg(h, data[:])
	// When: we deserialize the message
	// Then: we panic
	ParseNetlinkMsg(serializedData)
}

func TestParseNetlinkMessages(t *testing.T) {
	// Given: a list of serialized netlink messages
	h1 := CreateTestNlMsghdr()
	data1 := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	h1.Len = h1.Len + uint32(len(data1))
	h2 := CreateTestNlMsghdr()
	data2 := [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x01}
	h2.Len = h2.Len + uint32(len(data2))

	var nlmsgs []byte
	nlmsgs = append(nlmsgs, NewEncodedNetlinkMsg(h1, data1[:])...)
	nlmsgs = append(nlmsgs, NewEncodedNetlinkMsg(h2, data2[:])...)

	// When: parse these serialized data
	result := ParseNetlinkMsgs(nlmsgs)

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