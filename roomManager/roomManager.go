package roomManager

import(
	"github.com/iris-contrib/websocket"
	"github.com/jbrodriguez/mlog"
)

type Member struct{
	RoomID int64
	SendChan chan[] byte
	Conn *websocket.Conn
}


type RoomManager struct {
	CurRoomID int64
	Rooms *map[int64] *Room
}
func (r *RoomManager) isRoomNameExit(name string) (int64,bool) {
	for _,room  := range *r.Rooms{
		if room.Name == name {
			return room.ID,true
		}
	}
	return 0,false
}
func (r *RoomManager) CreateRoom(name string,m *Member) int64{

	if roomID,ok := r.isRoomNameExit(name) ; ok{
		r.JoinRoom(roomID,m)
		return roomID
	}else{
		mlog.Info("room not exit")
		curRoomID := r.CurRoomID
		m.RoomID = curRoomID
		members := []*Member{m}
		room := Room{ID:curRoomID,Name:name,Members: &members,Owner:m,FlvFirstPacket:nil}
		(*r.Rooms)[curRoomID] = &room
		r.CurRoomID = r.CurRoomID +1
		return r.CurRoomID
	}
}
func (r *RoomManager) DeleRoom( roomID int64){
	delete(*r.Rooms,roomID)
}

func (r *RoomManager) JoinRoom(roomID int64,m *Member){
	if roomTmp,ok := (*r.Rooms)[roomID] ; ok{
		m.RoomID = roomID
		roomTmp.AddMember(m)
	}else{
		mlog.Info("room not exit")
	}
}

func (r *RoomManager) MemberExit(roomID int64,m *Member){
	if room,ok := (*r.Rooms)[roomID];ok{
		room.DelMember(m)
	}
}
func NewRooms() *RoomManager{
	return &RoomManager{CurRoomID:0,Rooms:&map[int64]*Room{}}
}





type Room struct{
	ID int64
	Name string
	Owner *Member
	Members *[]*Member
	FlvFirstPacket  *[]byte
}


func (r *Room) AddMember(m *Member) {
	*r.Members = append(*r.Members,m)
}


func(r *Room) DelMember(m *Member){
	members := (*r).Members
	for i,member := range *members{
		if member != m{
			continue
		}
		newMembers := append((*members)[:i],(*members)[i+1:]...)
		(*r).Members =  &newMembers
	}
}

