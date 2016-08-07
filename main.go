

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


const (
	FLVPacketTypeSeq = 0
	FLVPacketTypeNlu = 1
    FLVPacketTypeEos = 2
)
func main() {
	mlog.StartEx(mlog.LevelWarn, "app.log", 5*1024*1024, 5)
	roomManeger := roomManager.NewRooms()
	r := router.New()
	r.Map(1000,func(m *roomManager.Member,p proto.Message) []byte{
		createRoom := p.(*myproto.CreateRoomTos)
		mlog.Info("create room :%s",createRoom.RoomName)
		roomID := roomManeger.CreateRoom(createRoom.RoomName,m)
		mlog.Info("room id %d",roomID)
		rData,_ := message.Marshal(&myproto.CreateRoomToc{RoomID:roomID})
		return rData
	})
	r.Map(1002,func(m *roomManager.Member,p proto.Message) []byte{
		joinRoom := p.(*myproto.JoinRoomTos)
		roomManeger.JoinRoom(joinRoom.RoomID,m)

		if data := (*roomManeger.Rooms)[m.RoomID].FlvFirstPacket;data !=nil {
			m.SendChan <- *data
			mlog.Warning("data %v", *data)
		}
		rData,_ := message.Marshal(&myproto.JoinRoomToc{ErrCode:joinRoom.RoomID})
		return rData
	})
	r.Map(1008,func(m *roomManager.Member,p proto.Message) []byte{
		var curRooms  []*myproto.Room
		mlog.Info("%v",*roomManeger.Rooms)
		for k,room := range *roomManeger.Rooms {
			roomID := k
			roomName := room.Name
			roomTmp := &myproto.Room{RoomID:roomID,RoomName:roomName}
			curRooms = append(curRooms,roomTmp)
		}
		rData,_ := message.Marshal(&myproto.GetRoomToc{Room:curRooms})
		mlog.Info("send %v",rData)
		return rData
	})

	r.LiveF= func(m *roomManager.Member,data *[]byte) []byte{
		members := (*roomManeger.Rooms)[m.RoomID].Members
		for _,roomMember := range *members {
			if roomMember != m {
				(*roomMember).SendChan <- *data
			}
		}
		
		// [Uint64-Len,Uint64-ProtoType,Uint8-streamType,Uint8-FlvHeader,Uint8-FlvPacketType,...]
		if (*data)[18] == FLVPacketTypeSeq {
			(*roomManeger.Rooms)[m.RoomID].FlvFirstPacket = data
		}
		return  []byte{}
	}

	upgrader := websocket.Custom(func( c *websocket.Conn){
		//mlog.Info("connect :%v",c)
		sendChan := make(chan []byte)
		m := roomManager.Member{SendChan:sendChan,Conn:c}
		go send(&m)
		if receive(&m,r){
			if room,ok := (*roomManeger.Rooms)[m.RoomID];ok{
				room.DelMember(&m)
				if len(*((*room).Members)) == 0 {
					roomManeger.DeleRoom(m.RoomID)
				}
			}
		}
	},100000,100000,true)
	wsHandlerFunc := func (ctx  *iris.Context){
		upgrader.Upgrade(ctx)
	}
	iris.Get("/ws", wsHandlerFunc)
	iris.Listen("0.0.0.0:8080")
}


func send(m *roomManager.Member){
	ticker := time.NewTicker( 10 * time.Second)
	defer func() {
		ticker.Stop()
		m.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-m.SendChan:
			if !ok {
				defer func() {
					if err := recover(); err != nil {
						ticker.Stop()
						m.Conn.Close()
					}
				}()
				m.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			m.Conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
			res, err := m.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			res.Write(msg)

			n := len(m.SendChan)
			for i := 0; i < n; i++ {
				res.Write(<-m.SendChan)
			}

			if err := res.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := m.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func receive(m *roomManager.Member,r *router.Router) bool {
	c := m.Conn
	s := m.SendChan
	defer func() {
		c.Close()
		close(m.SendChan)
	}()
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
				break
			}
		} else {
			mlog.Info("receive:%v",data)
			if needReply,sendData := r.DoRoute(m,&data);needReply{
				s <-sendData
			}
		}

	}

	return true
}