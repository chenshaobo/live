package rooms

import(
	"github.com/kataras/iris"
	"github.com/jbrodriguez/mlog"
)


type room struct{
	id int64
	name string
	member *[]iris.WebsocketConnection
}

type rooms struct {
	curRoomID int64
	rooms *map[string] *room
}

func (r *rooms) CreateRoom(name string,c iris.WebsocketConnection) int64{

	if roomTmp,ok := (*r.rooms)[name] ; ok{
		return roomTmp.id
	}else{
		mlog.Info("room not exit")
		curRoomID := r.curRoomID
		initMember := []iris.WebsocketConnection{c}
		(*r.rooms)[name] = &room{id:curRoomID,name:name,member: &initMember}
		r.curRoomID = r.curRoomID +1
		return r.curRoomID
	}
}

func (r *rooms) JoinRoom(name string,c iris.WebsocketConnection){
	if roomTmp,ok := (*r.rooms)[name] ; ok{
		(*roomTmp).member = append(*(*roomTmp).member,c)
	}else{
		mlog.Info("room not exit")
	}
}

func NewRooms() *rooms{
	return &rooms{curRoomID:0,rooms:&map[string] *room{}}
}



