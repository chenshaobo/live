package message

import(
	"github.com/golang/protobuf/proto"
	"github.com/jbrodriguez/mlog"
	"./myproto"
	"encoding/binary"
	"bytes"
)

const (
	MessageLenIndex = 8
	MessageTypeIndex = 16
)




func  Unmarshal(data *[]byte) (uint64,proto.Message) {
	mlog.Info("data:%v",*data)
	mlog.Info("data:%v",(*data)[:MessageLenIndex])
	messageLen := binary.BigEndian.Uint64((*data)[:MessageLenIndex])
	mlog.Info("message len :%d",messageLen)
	messageType := binary.BigEndian.Uint64((*data)[MessageLenIndex:MessageTypeIndex])
	pb :=(*myproto.Id2Name)[messageType]
	mlog.Info("message type :%d",messageType,"%+v",pb,"%v",(*data)[MessageTypeIndex:messageLen])
	err := proto.Unmarshal((*data)[MessageTypeIndex:messageLen],pb)
	if err != nil{
		panic(err)
	}
	return messageType,pb
}

func Marshal(message proto.Message) ([]byte,error){
	messageID := (*myproto.Name2IDStr)[proto.MessageName(message)]
	mlog.Info("message id :%v",messageID)
	tmpB := make([]byte,8)
	b := bytes.NewBuffer(make([]byte,0))


	if data,err := proto.Marshal(message);err != nil{
		return nil,err
	}else {

		binary.BigEndian.PutUint64(tmpB, uint64(8 + 8 + len(data)))
		mlog.Info("%v,len %d",tmpB,(8+8+len(data)))
		b.Write(tmpB)
		binary.BigEndian.PutUint64(tmpB,messageID)
		b.Write(tmpB)
		b.Write(data)
		return b.Bytes(),nil
	}
}