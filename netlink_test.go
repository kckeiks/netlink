package netlink

import (
	"bytes"
	"reflect"
	"testing"
	"golang.org/x/sys/unix"
)

func CreateTestNlMsghdr() unix.NlMsghdr {
	h := unix.NlMsghdr{}
	h.Len = uint32(unix.SizeofNlMsghdr)
	h.Type = uint16(2)
	h.Flags = uint16(5)
	h.Seq = uint32(6)
	h.Pid = uint32(11)
	return h
}

func TestNewSerializedNetlinkMsg(t *testing.T) {
	// Given: a NlMsghdr header and some data in bytes
	h := CreateTestNlMsghdr()
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	h.Len = h.Len + uint32(len(data))

	// When: we serialize the header and the data
	serializedData := NewSerializedNetlinkMsg(h, data[:])

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
	if testByteOrder.Uint32(data[:]) != testByteOrder.Uint32(serializedData[16:]) {
		t.Fatalf("NlMsghdr.Data = %d, expected %d", testByteOrder.Uint32(serializedData[16:]), data)
	}
	
	// Then: the serialized data has the correct number of bytes
	if uint32(len(serializedData)) != h.Len {
		t.Fatalf("Incorrect length len(serializedData)=%d, expected %d", len(serializedData), h.Len)
	}
}

func TestDeserializeNetlinkMessage(t *testing.T) {
	// Given: a serialized netlink message
	h := CreateTestNlMsghdr()
	data := [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	h.Len = h.Len + uint32(len(data))
	serializedData := NewSerializedNetlinkMsg(h, data[:])
	
	// When: we deserialize the message
	result, xdata := DeserializeNetlinkMsg(serializedData)

	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, h) {
		t.Fatalf("Given InetDiagReqV2 %+v and deserialized is %+v,", result, h)
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
	serializedData := NewSerializedNetlinkMsg(h, data)
	
	// When: we deserialize the message
	_, xdata := DeserializeNetlinkMsg(serializedData)

	// Then: nil is returned for the extra data
	if xdata != nil {
		t.Fatalf("Extra data=%d, expected nil.", xdata)
	}
}