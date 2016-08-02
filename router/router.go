package router

import (
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/live/message"
	"github.com/chenshaobo/live/roomManager"
)


type Router struct{
	routeMap *map[uint64] func(m *roomManager.Member,p proto.Message) []byte
}

func New() *Router {
	routeMap := make(map[uint64] (func( m *roomManager.Member,p proto.Message) []byte))
	return &Router{routeMap:&routeMap}
}
func(r *Router) Map(i uint64,f func( m *roomManager.Member,p proto.Message) []byte){
	(*r.routeMap)[i] = f
}
func(r *Router) GetRouteFun(i uint64) func(m *roomManager.Member,p proto.Message) []byte{
	return (*r.routeMap)[i]
}
func(r *Router) DoRoute(m *roomManager.Member, data *[]byte) []byte {
	messageType,protoMsg := message.Unmarshal(data)
	routerFun := (*r.routeMap)[messageType]
	return  routerFun(m,protoMsg)
}
