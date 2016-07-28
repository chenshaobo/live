

package main
import (
	"github.com/kataras/iris"
	"github.com/jbrodriguez/mlog"
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/live/router"
	"github.com/chenshaobo/live/myproto"
	"github.com/chenshaobo/live/message"
	"github.com/chenshaobo/live/room"
)



func main() {
	mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)
	iris.Config.Websocket.Endpoint = "/ws"
	roomManeger := rooms.NewRooms()
	r := router.New()
	r.Map(1000,func(c iris.WebsocketConnection,p proto.Message) []byte{
		createRoom := p.(*myproto.CreateRoomTos)
		mlog.Info("create room :%s",createRoom.RoomName)
		roomID := "1000"
		roomManeger.CreateRoom(createRoom.RoomName,c)
		rData,_ := message.Marshal(&myproto.CreateRoomToc{RoomID:&roomID})
		return rData
	})
	iris.Websocket.OnConnection(func(c iris.WebsocketConnection){
		mlog.Info("Connect websocket")
		c.OnMessage(func(message []byte){
			mlog.Info("receive %v",message)
			sendMsg := r.DoRoute(c,&message)

			c.EmitMessage(sendMsg)
		})
	})

	iris.Listen(":8080")
}
