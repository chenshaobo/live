package main

import (
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/live/myproto"
	"github.com/chenshaobo/live/message"
	"golang.org/x/net/websocket"
	"github.com/golang/protobuf/proto"
)

func main() {
	mlog.StartEx(mlog.LevelInfo, "app.log", 5 * 1024 * 1024, 5)
	ws, err := websocket.Dial("ws://127.0.0.1:8080/ws", "", "http://127.0.0.1:8080")
	if err != nil {
		mlog.Error(err)
		panic(err)
	}
	test(1000,ws)
	test(1008,ws)
	for {

	}
}

func test(msgType uint64,ws * websocket.Conn){
	var msg = new([]byte)
	switch msgType {
	case 1000:
		sendData(&myproto.CreateRoomTos{RoomName:"dddd"},ws)
	case 1008:
		sendData(&myproto.GetRoomsTos{Id:1000},ws)

	}
	err := websocket.Message.Receive(ws, msg)
	isErr(err)
	_, dataProto := message.Unmarshal(msg)
	mlog.Info("%+v", dataProto)




}

func sendData(p proto.Message,ws *websocket.Conn){
	data,err := message.Marshal(p)
	mlog.Info("data:%v", data)
	isErr(err)
	ws.PayloadType = websocket.BinaryFrame
	_, err1 := ws.Write(data)
	isErr(err1)
}
func isErr(err error) {
	if err != nil {
		mlog.Error(err)
		panic(err)
	}
}