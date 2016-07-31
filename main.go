

package main
import (
	websocket "github.com/iris-contrib/websocket"
	//"github.com/kataras/iris"
	//"github.com/valyala/fasthttp"
	"github.com/jbrodriguez/mlog"
	"github.com/kataras/iris"
	"time"
	"github.com/chenshaobo/live/roomManager"
	"github.com/chenshaobo/live/router"
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/live/myproto"
	"github.com/chenshaobo/live/message"
)



func main() {
	mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)
	roomManeger := roomManager.NewRooms()
	r := router.New()
	r.Map(1000,func(c *websocket.Conn,p proto.Message) []byte{
		createRoom := p.(*myproto.CreateRoomTos)
		mlog.Info("create room :%s",createRoom.RoomName)
		roomID := roomManeger.CreateRoom(createRoom.RoomName,c)
		mlog.Info("room id %d",roomID)
		rData,_ := message.Marshal(&myproto.CreateRoomToc{RoomID:roomID})
		return rData
	})

	r.Map(1008,func(c *websocket.Conn,p proto.Message) []byte{
		var curRooms  []*myproto.Room
		mlog.Info("%v",*roomManeger.Rooms)
		for k,room := range *roomManeger.Rooms {
			roomName := k
			roomID := room.ID
			roomTmp := &myproto.Room{RoomID:roomID,RoomName:roomName}
			curRooms = append(curRooms,roomTmp)
		}
		rData,_ := message.Marshal(&myproto.GetRoomToc{Room:curRooms})
		mlog.Info("send %v",rData)
		return rData
	})

	upgrader := websocket.Custom(func( c *websocket.Conn){
		//mlog.Info("connect :%v",c)
		sendChan := make(chan []byte)
		go send(sendChan,c)
		receive(sendChan,c,r)
	},100000,100000,true)
	wsHandlerFunc := func (ctx  *iris.Context){
		upgrader.Upgrade(ctx)
	}
	iris.Get("/ws", wsHandlerFunc)
	iris.Listen("0.0.0.0:8080")
}


func send(s chan []byte,c *websocket.Conn){
	ticker := time.NewTicker( 10 * time.Second)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-s:
			if !ok {
				defer func() {
					if err := recover(); err != nil {
						ticker.Stop()
						c.Close()
					}
				}()
				c.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.SetWriteDeadline(time.Now().Add(60 * time.Second))
			res, err := c.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			res.Write(msg)

			n := len(s)
			for i := 0; i < n; i++ {
				res.Write(<-s)
			}

			if err := res.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				mlog.Error(err)
				return
			}
		}
	}
}

func receive(s chan []byte,c *websocket.Conn,r *router.Router) {
	c.SetReadLimit(1024*1024)
	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetPongHandler(func(s string) error {
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		if _, data, err := c.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				mlog.Error(err)
				c.Close()
				break
			}
		} else {
			mlog.Info("receive:%v",data)
			sendData := r.DoRoute(c,&data)
			s<- sendData
		}

	}
}