

package main
import (
	"github.com/kataras/iris"
	"github.com/jbrodriguez/mlog"
	"github.com/golang/protobuf/proto"
	"./router"
	"./myproto"
	"./message"
)



func main() {
	mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)
	iris.Config.Websocket.Endpoint = "/ws"
	r := router.New()
	r.Map(1000,func(p proto.Message) []byte{
		createRoom := p.(*myproto.CreateRoomTos)
		mlog.Info("create room :%s",createRoom.RoomName)
		roomID := "1000"
		rData,_ := message.Marshal(&myproto.CreateRoomToc{RoomID:&roomID})
		return rData
	})
	iris.Websocket.OnConnection(func(c iris.WebsocketConnection){
		mlog.Info("Connect websocket")
		c.OnMessage(func(message []byte){
			mlog.Info("receive %v",message)
			sendMsg := r.DoRoute(&message)

			c.EmitMessage(sendMsg)
		})
	})

	iris.Listen(":8080")
}
