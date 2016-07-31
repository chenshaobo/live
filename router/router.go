package router

import (
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/live/message"
	"github.com/iris-contrib/websocket"
)


type Router struct{
	routeMap *map[uint64] func( *websocket.Conn,proto.Message) []byte
}

func New() *Router {
	routeMap := make(map[uint64] (func( *websocket.Conn,proto.Message) []byte))
	return &Router{routeMap:&routeMap}
}
func(r *Router) Map(i uint64,f func( *websocket.Conn,proto.Message) []byte){
	(*r.routeMap)[i] = f
}
func(r *Router) GetRouteFun(i uint64) func(*websocket.Conn,proto.Message) []byte{
	return (*r.routeMap)[i]
}
func(r *Router) DoRoute(c *websocket.Conn, data *[]byte) []byte {
	messageType,protoMsg := message.Unmarshal(data)
	routerFun := (*r.routeMap)[messageType]
	return  routerFun(c,protoMsg)
}
