package myproto

import (
	"github.com/golang/protobuf/proto"
)

var Id2Name *map[uint64]proto.Message = makeID2Name()
var Name2IDStr *map[string]uint64 = makeName2ID()

func init() {
	Id2Name = makeID2Name()
	Name2IDStr = makeName2ID()
}
func makeID2Name() *map[uint64]proto.Message {
	return &map[uint64]proto.Message{
		1000 :&CreateRoomTos{},
		1001 : &CreateRoomToc{},
		1002 : &JoinRoomTos{},
		1003 : &JoinRoomToc{},
		1004 : &LeaveRoomTos{},
		1005 : &LeaveRoomToc{},
		1006 : &LiveTos{},
		1007 : &LiveToc{},
		1008 : &GetRoomsTos{},
		1009 : &GetRoomToc{}, }
}

func makeName2ID() *map[string]uint64 {
	return &map[string]uint64{
		proto.MessageName(&CreateRoomTos{}) : 1000,
		proto.MessageName(&CreateRoomToc{}) : 1001,
		proto.MessageName(&JoinRoomTos{}) : 1002,
		proto.MessageName(&JoinRoomToc{}) : 1003,
		proto.MessageName(&LeaveRoomTos{}) : 1004,
		proto.MessageName(&LeaveRoomToc{}) : 1005,
		proto.MessageName(&LiveTos{}) : 1006,
		proto.MessageName(&LiveToc{}) : 1007,
		proto.MessageName(&GetRoomsTos{}) : 1008,
		proto.MessageName(&GetRoomToc{}) : 1009, }
}
