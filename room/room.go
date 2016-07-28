package rooms

import(
	"github.com/kataras/iris"
)

type rooms struct {
	rooms map[string] room
}

func (r *rooms) CreateRoom(name string,c iris.WebsocketConnection){

}
func NewRooms() *rooms{
	return &rooms{}
}



type room struct{
	id uint64
	name string
	member []iris.WebsocketConnection
}