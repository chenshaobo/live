package roomManager

import(
	"github.com/iris-contrib/websocket"
	"github.com/jbrodriguez/mlog"
)


type RoomManager struct {
	CurRoomID int64
	Rooms *map[string] *room
}

func (r *RoomManager) CreateRoom(name string,c *websocket.Conn) int64{

	if roomTmp,ok := (*r.Rooms)[name] ; ok{
		return roomTmp.ID
	}else{
		mlog.Info("room not exit")
		curRoomID := r.CurRoomID
		initMember := []*websocket.Conn{c}
		(*r.Rooms)[name] = &room{ID:curRoomID,Name:name,member: &initMember}
		r.CurRoomID = r.CurRoomID +1
		return r.CurRoomID
	}
}

func (r *RoomManager) JoinRoom(name string,c *websocket.Conn){
	if roomTmp,ok := (*r.Rooms)[name] ; ok{
		roomTmp.AddMember(c)
	}else{
		mlog.Info("room not exit")
	}
}

func NewRooms() *RoomManager{
	return &RoomManager{CurRoomID:0,Rooms:&map[string] *room{}}
}





type room struct{
	ID int64
	Name string
	member *[]*websocket.Conn
}


func (r *room) AddMember(c *websocket.Conn) {
	*r.member = append(*r.member,c)
}

