package rooms

import(
	"github.com/kataras/iris"
	"github.com/jbrodriguez/mlog"
)


type rooms struct {
	curRoomID int64
	Rooms *map[string] *room
}

func (r *rooms) CreateRoom(name string,c iris.WebsocketConnection) int64{

	if roomTmp,ok := (*r.Rooms)[name] ; ok{
		return roomTmp.ID
	}else{
		mlog.Info("room not exit")
		curRoomID := r.curRoomID
		initMember := []iris.WebsocketConnection{c}
		(*r.Rooms)[name] = &room{ID:curRoomID,Name:name,member: &initMember}
		r.curRoomID = r.curRoomID +1
		return r.curRoomID
	}
}

func (r *rooms) JoinRoom(name string,c iris.WebsocketConnection){
	if roomTmp,ok := (*r.Rooms)[name] ; ok{
		roomTmp.AddMember(c)
	}else{
		mlog.Info("room not exit")
	}
}

func NewRooms() *rooms{
	return &rooms{curRoomID:0,Rooms:&map[string] *room{}}
}





type room struct{
	ID int64
	Name string
	member *[]iris.WebsocketConnection
}


func (r *room) AddMember(c iris.WebsocketConnection) {
	*r.member = append(*r.member,c)
}

