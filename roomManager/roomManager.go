package roomManager

import(
	"github.com/iris-contrib/websocket"
	"github.com/jbrodriguez/mlog"
)

type Member struct{
	SendChan chan[] byte
	Conn *websocket.Conn
}


type RoomManager struct {
	CurRoomID int64
	Rooms *map[string] *Room
	Roomers *map[Member] *Room
}

func (r *RoomManager) CreateRoom(name string,m *Member) int64{

	if roomTmp,ok := (*r.Rooms)[name] ; ok{
		r.JoinRoom(name,m)
		return roomTmp.ID
	}else{
		mlog.Info("room not exit")
		curRoomID := r.CurRoomID
		initMember := []*Member{m}
		room := Room{ID:curRoomID,Name:name,Members: &initMember}
		mlog.Info("%v",m)
		(*r.Rooms)[name] = &room
		mlog.Info("%v, %v",m,*r.Roomers)
		(*r.Roomers)[*m] = &room
		r.CurRoomID = r.CurRoomID +1
		return r.CurRoomID
	}
}

func (r *RoomManager) JoinRoom(name string,m *Member){
	if roomTmp,ok := (*r.Rooms)[name] ; ok{
		roomTmp.AddMember(m)
	}else{
		mlog.Info("room not exit")
	}
}

func NewRooms() *RoomManager{
	return &RoomManager{CurRoomID:0,Rooms:&map[string] *Room{},Roomers:&map[Member] *Room{}}
}





type Room struct{
	ID int64
	Name string
	Members *[]*Member
}


func (r *Room) AddMember(m *Member) {
	*r.Members = append(*r.Members,m)
}



