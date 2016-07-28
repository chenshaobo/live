package main

import (
"regexp"
"io/ioutil"
"os"

)

func main() {
	// Compile the expression once, usually at init time.
	// Use raw strings to avoid having to quote the backslashes.
	//message[ \t]*[ \t\n]*(\/\/<(.*)>)*[ \t\n]*[ \t\n]*{[ \t\n]*[^}]*}
	var validID = regexp.MustCompile("message[ \\t]+([\\w]+)[ \\t\\r\\n]*({\\/\\/<([\\w]+)>)*[ \\t\\n]*[ \\t\\n]*([ \\t\\r\\n]+[^}]*)}")
	d,_ :=ioutil.ReadFile("../myproto.proto")
	reg := validID.FindAllSubmatch(d,-1)

	id2NameStr :="package myproto \n import (\n \"github.com/golang/protobuf/proto\") \n " +
		         " var Id2Name *map[uint64]proto.Message = makeID2Name() \n var Name2IDStr *map[string]uint64 = makeName2ID()\n" +
				"func init(){ \n Id2Name = makeID2Name()\n Name2IDStr = makeName2ID()\n}\n"+
				"func makeID2Name() *map[uint64]proto.Message { \n return &map[uint64]proto.Message{"
	name2IDStr := "func makeName2ID() *map[string]uint64 { \n return &map[string]uint64{"
	strTmp1 :="\n"
	strTmp2 :="\n"
	for _,match := range reg {

		name := match[1]
		id := match[3]
		if strTmp1 != "\n" {
			strTmp1 = strTmp1 +",\n"  +  string(id) + " : &" +string(name) + "{}"
			strTmp2 = strTmp2 + ",\n"  + "proto.MessageName(&" + string(name) + "{} )" + " : " + string(id)
		} else {
			strTmp1 = strTmp1 + string(id) + " :&" +string(name) + "{}"
			strTmp2 =  strTmp2 + "proto.MessageName(&" + string(name) + "{} )" + " : " + string(id)
		}
	}
	id2NameStr = id2NameStr + strTmp1 + ",}}\n"
	name2IDStr = name2IDStr + strTmp2 + ",}}\n"

	f, err := os.OpenFile("../protoMap.go", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(id2NameStr + "\n" +name2IDStr); err != nil {
		panic(err)
	}

}