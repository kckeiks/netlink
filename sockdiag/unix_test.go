package sockdiag

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/kckeiks/netlink/internal/testutils"
)

func newTestUnixDiagReq() UnixDiagReq {
	req := UnixDiagReq{}
	req.Family = 1
	req.Protocol = 2
	req.Pad = 3
	req.States = 4
	req.Inode = 5
	req.Show = 6
	req.Cookie = [2]uint32{7, 8}
	return req
}

func newTestUnixDiagMsg() UnixDiagMsg {
	idm := UnixDiagMsg{}
	idm.Family = 1
	idm.Type = 2
	idm.State = 3
	idm.Pad = 4
	idm.Inode = 5
	return idm
}

func deserializeUnixDiagReq(data []byte) UnixDiagReq {
	b := bytes.NewBuffer(data)
	req := UnixDiagReq{}
	err := binary.Read(b, testutils.TestByteOrder, &req)
	if err != nil {
		panic("Error: Could not deserialize UnixDiagReq.")
	}
	return req
}

func TestSerializeUnixDiagReq(t *testing.T) {
	// Given: a unix_diag_req header
	req := newTestUnixDiagReq()
	// When: we serialize the header
	serializedData, err := SerializeUnixDiagReq(req)
	if err != nil {
		t.Fatalf("got an unexpected error %v.", err)
	}
	// Then: it's serialized with the correct data
	if req.Family != serializedData[0] {
		t.Fatalf("UnixDiagReq.Family = %d, expected %d", serializedData[0], req.Family)
	}
	if req.Protocol != serializedData[1] {
		t.Fatalf("UnixDiagReq.Protocol = %d, expected %d", serializedData[1], req.Protocol)
	}
	if req.Pad != testutils.TestByteOrder.Uint16(serializedData[2:4]) {
		t.Fatalf("UnixDiagReq.Pad = %d, expected %d", serializedData[2:4], req.Pad)
	}
	if req.States != testutils.TestByteOrder.Uint32(serializedData[4:8]) {
		t.Fatalf("UnixDiagReq.States = %d, expected %d", testutils.TestByteOrder.Uint32(serializedData[4:8]), req.States)
	}
	if req.Inode != testutils.TestByteOrder.Uint32(serializedData[8:12]) {
		t.Fatalf("UnixDiagReq.Inode = %d, expected %d", testutils.TestByteOrder.Uint32(serializedData[8:12]), req.Inode)
	}
	if req.Show != testutils.TestByteOrder.Uint32(serializedData[12:16]) {
		t.Fatalf("UnixDiagReq.Show = %d, expected %d", testutils.TestByteOrder.Uint32(serializedData[12:16]), req.Show)
	}
	if req.Cookie[0] != testutils.TestByteOrder.Uint32(serializedData[16:20]) {
		t.Fatalf("UnixDiagReq.Cookie[0] = %d, expected %d", req.Cookie[0], testutils.TestByteOrder.Uint32(serializedData[16:20]))
	}
	if req.Cookie[1] != testutils.TestByteOrder.Uint32(serializedData[20:24]) {
		t.Fatalf("UnixDiagReq.Cookie[1] = %d, expected %d", req.Cookie[1], testutils.TestByteOrder.Uint32(serializedData[20:24]))
	}
}

func TestDeserializeUnixDiagMsg(t *testing.T) {
	// Given: a serialized InetDiagMsg
	msg := newTestUnixDiagMsg()
	serializedData := bytes.NewBuffer(make([]byte, NlUnixDiagMsgLen))
	serializedData.Reset()
	binary.Write(serializedData, testutils.TestByteOrder, &msg)
	// When: we deserialize
	result, err := DeserializeUnixDiagMsg(serializedData.Bytes())
	// Then: there is no error
	if err != nil {
		t.Fatalf("got an unexpected error %v.", err)
	}
	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, &msg) {
		t.Fatalf("Given UnixDiagMsg %+v but expected %+v,", result, msg)
	}
}

func TestDeserializeUnixDiagReq(t *testing.T) {
	// Given: a unix_diag_req header
	req := newTestUnixDiagReq()
	serializedData, err := SerializeUnixDiagReq(req)
	if err != nil {
		t.Fatalf("got an unexpected error %v.", err)
	}
	// When: we deserialize
	result := deserializeUnixDiagReq(serializedData)
	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, req) {
		t.Fatalf("Given UnixDiagReq %+v and deserialized is %+v,", req, result)
	}
}

func TestNewUnixNetlinkMsg(t *testing.T) {
	// Given: a NlMsghdr header and some data in bytes
	h := testutils.NewTestNlMsghdr()
	unixHeader := newTestUnixDiagReq()
	h.Len = NlUnixDiagReqMsgLen
	// When: we serialize the header and the data
	serializedData, err := NewUnixNetlinkMsg(h, unixHeader)
	// Then: there is no error
	if err != nil {
		t.Fatalf("got an unexpected error %v.", err)
	}
	// Then: the message was serialized with the correct data
	if h.Len != testutils.TestByteOrder.Uint32(serializedData[:4]) {
		t.Fatalf("NlMsghdr.Length = %d, expected %d", testutils.TestByteOrder.Uint32(serializedData[:4]), h.Len)
	}
	if h.Type != testutils.TestByteOrder.Uint16(serializedData[4:6]) {
		t.Fatalf("NlMsghdr.Type = %d, expected %d", testutils.TestByteOrder.Uint16(serializedData[4:6]), h.Type)
	}
	if h.Flags != testutils.TestByteOrder.Uint16(serializedData[6:8]) {
		t.Fatalf("NlMsghdr.Flags = %d, expected %d", testutils.TestByteOrder.Uint16(serializedData[6:8]), h.Flags)
	}
	if h.Seq != testutils.TestByteOrder.Uint32(serializedData[8:12]) {
		t.Fatalf("NlMsghdr.Seq = %d, expected %d", testutils.TestByteOrder.Uint32(serializedData[8:12]), h.Seq)
	}
	if h.Pid != testutils.TestByteOrder.Uint32(serializedData[12:16]) {
		t.Fatalf("NlMsghdr.Pid = %d, expected %d", testutils.TestByteOrder.Uint32(serializedData[12:16]), h.Pid)
	}
	deserializedUnixHeader := deserializeUnixDiagReq(serializedData[16:])
	// Then: the deserialized inetHeader that we get has the same values as the initial inetHeader
	if !reflect.DeepEqual(deserializedUnixHeader, unixHeader) {
		t.Fatalf("Given UnixDiagReq %+v and deserialized is %+v,", unixHeader, deserializedUnixHeader)
	}
	// Then: the serialized data has the correct number of bytes
	if uint32(len(serializedData)) != h.Len {
		t.Fatalf("Incorrect length len(serializedData)=%d, expected %d", len(serializedData), h.Len)
	}
}
