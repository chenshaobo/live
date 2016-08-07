package router

import (
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/live/message"
	"github.com/chenshaobo/live/roomManager"
	"encoding/binary"
)


type Router struct{
	routeMap *map[uint64] func(m *roomManager.Member,p proto.Message) []byte
	LiveF   func(m *roomManager.Member,data *[]byte)[]byte
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
func(r *Router) DoRoute(m *roomManager.Member, data *[]byte) (bool,[]byte) {
	protoType  := binary.BigEndian.Uint64((*data)[8:16])
    if protoType == 255 {
		return false,r.LiveF(m,data)
	}else {
		messageType, protoMsg := message.Unmarshal(data)
		routerFun := (*r.routeMap)[messageType]
		return true,routerFun(m, protoMsg)
	}
}
