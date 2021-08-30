package netlink

import (
	"testing"
	"bytes"
	"reflect"
	"encoding/binary"
)

func CreateTestInetDiagReqV2() InetDiagReqV2 {
	req := InetDiagReqV2{}
	req.Family = 1
	req.Protocol = 2
	req.Ext = 3
	req.Pad = 4
	req.States = 5
	idsi := createTestInetDiagSockID()
	req.ID = idsi
	return req
}

func createTestInetDiagMsg() InetDiagMsg {
	idm := InetDiagMsg{}
	idm.Family = 1
	idm.State = 2
	idm.Timer = 3
	idm.Retrans = 4
	idm.ID = createTestInetDiagSockID()
	idm.Expires = 5
	idm.RQueue = 6
	idm.WQueue = 7
	idm.UID = 8
	idm.Inode = 9
	return idm
}

func createTestInetDiagSockID() InetDiagSockID {
	idsi := InetDiagSockID{}
	idsi.SPort = [2]byte{0x10, 0x20}
	idsi.DPort = [2]byte{0x30, 0x40}
	idsi.Src = [16]byte{
		0x10, 0x20, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10,
		0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x90,
	}
	idsi.Dst = [16]byte{
		0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10,
		0x10, 0x20, 0x10, 0x10, 0x10, 0x10, 0x10, 0x70,
	}
	idsi.If = 6
	idsi.Cookie = [2]uint32{7, 8}
	return idsi
}

func TestSerializeInetDiagReqV2(t *testing.T) {
	// Given: a inet_diag_req_v2 header
	req := CreateTestInetDiagReqV2()
	
	// When: we serialize the header
	serializedData := SerializeInetDiagReqV2(req)

	// Then: it's serialized with the correct data
	if req.Family != serializedData[0] {
		t.Fatalf("InetDiagReqV2.Family = %d, expected %d", serializedData[0], req.Family)
	}
	if req.Protocol != serializedData[1] {
		t.Fatalf("InetDiagReqV2.Protocol = %d, expected %d", serializedData[1], req.Protocol)
	}
	if req.Ext != serializedData[2] {
		t.Fatalf("InetDiagReqV2.Ext = %d, expected %d", serializedData[2], req.Ext)
	}
	if req.Pad != serializedData[3] {
		t.Fatalf("InetDiagReqV2.Pad = %d, expected %d", serializedData[3], req.Pad)
	}
	if req.States != testByteOrder.Uint32(serializedData[4:8]) {
		t.Fatalf("InetDiagReqV2.States = %d, expected %d", testByteOrder.Uint32(serializedData[4:8]), req.States)
	}
	if bytes.Compare(req.ID.SPort[:], serializedData[8:10]) != 0 {
		t.Fatalf("InetDiagReqV2.ID.SPort = %d, expected %d", req.ID.SPort, serializedData[8:10])
	}
	if bytes.Compare(req.ID.DPort[:], serializedData[10:12]) != 0 {
		t.Fatalf("InetDiagReqV2.ID.DPort = %d, expected %d", req.ID.DPort, serializedData[10:12])
	}
	if bytes.Compare(req.ID.Src[:], serializedData[12:28]) != 0 {
		t.Fatalf("InetDiagReqV2.ID.Src = %d, expected %d", req.ID.Src, serializedData[12:28])
	}
	if bytes.Compare(req.ID.Dst[:], serializedData[28:44]) != 0 {
		t.Fatalf("InetDiagReqV2.ID.Dst = %d, expected %d", req.ID.Dst, serializedData[28:44])
	}
	if req.ID.If != testByteOrder.Uint32(serializedData[44:48]) {
		t.Fatalf("InetDiagReqV2.ID.If = %d, expected %d", req.ID.If, serializedData[44:48])
	}
	if req.ID.Cookie[0] != testByteOrder.Uint32(serializedData[48:52]) {
		t.Fatalf("InetDiagReqV2.ID.Cookie[0] = %d, expected %d", req.ID.Cookie[0], testByteOrder.Uint32(serializedData[48:52]))
	}
	if req.ID.Cookie[1] != testByteOrder.Uint32(serializedData[52:56]) {
		t.Fatalf("InetDiagReqV2.ID.Cookie[1] = %d, expected %d", req.ID.Cookie[1], testByteOrder.Uint32(serializedData[52:56]))
	}
}

func TestDeserializeInetDiagReqV2(t *testing.T) {
	// Given: a inet_diag_req_v2 header
	req := CreateTestInetDiagReqV2()
	serializedData := SerializeInetDiagReqV2(req)

	// When: we deserialize
	result := DeserializeInetDiagReqV2(serializedData)

	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, req) {
		t.Fatalf("Given InetDiagReqV2 %+v and deserialized is %+v,", req, result)
	}
}

func TestParseInetDiagMsg(t *testing.T) {
	// Given: a serialized InetDiagMsg
	msg := createTestInetDiagMsg()
	serializedData := bytes.NewBuffer(make([]byte, sizeOfInetDiagMsg))
	serializedData.Reset()
	binary.Write(serializedData, testByteOrder, &msg)

	// When: we deserialize
	result := ParseInetDiagMsg(serializedData.Bytes())

	// Then: the struct that we get has the same values as the initial struct
	if !reflect.DeepEqual(result, msg) {
		t.Fatalf("Given InetDiagMsg %+v but expected %+v,", result, msg)
	}
}