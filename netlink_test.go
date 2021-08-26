package netlink

import (
	"bytes"
	"reflect"
	"testing"
	"golang.org/x/sys/unix"
)

func CreateTestNetlinkMessage() NetlinkMessage {
	m := NetlinkMessage{}
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	m.Data = data[:]
	
	header := unix.NlMsghdr{}
	header.Len = uint32(unix.SizeofNlMsghdr + len(data))
	header.Type = uint16(2)
	header.Flags = uint16(5)
	header.Seq = uint32(6)
	header.Pid = uint32(11)
	m.Header = header
	return m
}

func TestSerializeNetlinkMessage(t *testing.T) {
	// Given: a netlink message with random values
	m := CreateTestNetlinkMessage()
	
	// When: we serialize the message
	serializedData := SerializeNetlinkMessage(m)

	// Then: the message was serialized with the correct data
	if m.Header.Len != testByteOrder.Uint32(serializedData[:4]) {
		t.Fatalf("NlMsghdr.Length = %d, expected %d", testByteOrder.Uint32(serializedData[:4]), m.Header.Len)
	}
	if m.Header.Type != testByteOrder.Uint16(serializedData[4:6]) {
		t.Fatalf("NlMsghdr.Type = %d, expected %d", testByteOrder.Uint16(serializedData[4:6]), m.Header.Type)
	}
	if m.Header.Flags != testByteOrder.Uint16(serializedData[6:8]) {
		t.Fatalf("NlMsghdr.Flags = %d, expected %d", testByteOrder.Uint16(serializedData[6:8]), m.Header.Flags)
	}
	if m.Header.Seq != testByteOrder.Uint32(serializedData[8:12]) {
		t.Fatalf("NlMsghdr.Seq = %d, expected %d", testByteOrder.Uint32(serializedData[8:12]), m.Header.Seq)
	}
	if m.Header.Pid != testByteOrder.Uint32(serializedData[12:16]) {
		t.Fatalf("NlMsghdr.Pid = %d, expected %d", testByteOrder.Uint32(serializedData[12:16]), m.Header.Pid)
	}
	if testByteOrder.Uint32(m.Data) != testByteOrder.Uint32(serializedData[16:]) {
		t.Fatalf("NlMsghdr.Data = %d, expected %d", testByteOrder.Uint32(serializedData[16:]), m.Data)
	}
	
	// Then: the serialized data has the correct number of bytes
	if uint32(len(serializedData)) != m.Header.Len {
		t.Fatalf("Incorrect length len(serializedData)=%d, expected %d", len(serializedData), m.Header.Len)
	}
}

func TestDeserializeNetlinkMessage(t *testing.T) {
	// Given: a serialized netlink message
	m := CreateTestNetlinkMessage()
	serializedData := SerializeNetlinkMessage(m)
	
	// When: we deserialize the message
	result, xdata := deserializeNetlinkMessage(serializedData)

	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, m.Header) {
		t.Fatalf("Given InetDiagReqV2 %+v and deserialized is %+v,", result, m.Header)
	}
	// Then: the extra data was returned
	if bytes.Compare(xdata, m.Data) != 0 {
		t.Fatalf("Extra data=%d, expected %d", xdata, m.Data)
	}
}

func TestDeserializeNetlinkMessageWithOutData(t *testing.T) {
	// Given: a serialized netlink message without extra data
	m := CreateTestNetlinkMessage()
	m.Data = []byte{} 
	m.Header.Len = uint32(unix.SizeofNlMsghdr)
	serializedData := SerializeNetlinkMessage(m)
	
	// When: we deserialize the message
	_, xdata := deserializeNetlinkMessage(serializedData)

	// Then: nil is returned for the extra data
	if xdata != nil {
		t.Fatalf("Extra data=%d, expected %d", xdata, m.Data)
	}
}