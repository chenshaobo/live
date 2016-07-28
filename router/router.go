package router

import (
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/live/message"
	"github.com/kataras/iris"
)


type Router struct{
	routeMap *map[uint64] func( iris.WebsocketConnection,proto.Message) []byte
}

func New() *Router {
	routeMap := make(map[uint64] (func( iris.WebsocketConnection,proto.Message) []byte))
	return &Router{routeMap:&routeMap}
}
func(r *Router) Map(i uint64,f func( iris.WebsocketConnection,proto.Message) []byte){
	(*r.routeMap)[i] = f
}
func(r *Router) GetRouteFun(i uint64) func(iris.WebsocketConnection,proto.Message) []byte{
	return (*r.routeMap)[i]
}
func(r *Router) DoRoute(c iris.WebsocketConnection, data *[]byte) []byte {
	messageType,protoMsg := message.Unmarshal(data)
	routerFun := (*r.routeMap)[messageType]
	return  routerFun(c,protoMsg)
}


//func getMessageType(data *[]byte) uint64{
//	return binary.BigEndian.Uint64((*data)[MessageLenIndex:MessageTypeIndex])
//}
//
//func getMessageLen(data *[]byte) uint64{
//	return binary.BigEndian.Uint64((*data)[:MessageLenIndex])
//}