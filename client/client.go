package main

import (
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/live/myproto"
	"github.com/chenshaobo/live/message"
	"golang.org/x/net/websocket"
	"os"
)

func main() {
	roomName := os.Args[1]
	mlog.StartEx(mlog.LevelInfo, "app.log", 5 * 1024 * 1024, 5)
	ws, err := websocket.Dial("ws://127.0.0.1:8080/ws", "", "http://127.0.0.1:8080")
	if err != nil {
		mlog.Error(err)
		panic(err)
	}
	var msg = new([]byte)

	data, err := message.Marshal(&myproto.CreateRoomTos{RoomName:&roomName})
	mlog.Info("data:%v", data)
	isErr(err)
	ws.PayloadType = websocket.BinaryFrame
	_, err1 := ws.Write(data)
	isErr(err1)

	err = websocket.Message.Receive(ws, msg)
	isErr(err)
	_, dataProto := message.Unmarshal(msg)
	mlog.Info("%+v", dataProto)
	mlog.Info("receive :%v", dataProto.(*myproto.CreateRoomToc).GetRoomID())

	data,err = message.Marshal(&myproto.GetRoomsTos{})
	mlog.Info("data:%v", data)
	isErr(err)
	ws.PayloadType = websocket.BinaryFrame
	_, err1 = ws.Write(data)
	isErr(err1)

	err = websocket.Message.Receive(ws, msg)
	isErr(err)
	_, dataProto = message.Unmarshal(msg)
	mlog.Info("%+v", dataProto)
	mlog.Info("receive :%v", *dataProto.(*myproto.GetRoomToc).GetRoom()[0].RoomName)
}

func isErr(err error) {
	if err != nil {
		mlog.Error(err)
		panic(err)
	}
}