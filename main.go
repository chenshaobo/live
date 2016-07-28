//package main
//
//import (
//	"github.com/jbrodriguez/mlog"
//	//	"github.com/labstack/echo"
//	//	"github.com/labstack/echo/engine/standard"
//	//	"net/http"
//	"github.com/streadway/amqp"
//	"os/exec"
//	"runtime"
//	"work/rabbitmq"
//	//	"work/proto"
//)


//func main() {
	//mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)
	//m := martini.Classic()
	//m.Get("/", func() string {
	//	return "Hello world!"
	//})
	//m.Get("/hello", func() string {
	//	return "Hello world!"
	//})
	//go m.Run()
	//list := []int{1, 2, 3, 4, 3}
	//mlog.Info("list 2:~s", list[:3])
	//mlog.Info("runcpu %d", runtime.NumCPU())
	//mq := rabbitmq.RabbitMQ{"amqp://guest:guest@172.16.2.109:5672/", "fanout_logger", "fanout", "fanout_logger_queue11", "fanout_logger_rk"}
	//go mq.ConnectRabbitmq(func(deliver amqp.Delivery) error {
	//	mlog.Info("%s", deliver.Body)
	//	return nil
	//})
	//e := echo.New()
	//e.Post("/api/image/convert/webp", func(c echo.Context) error {
	//	source := new(proto.ConvertImageJson)
	//	if error := c.Bind(source);error != nil{
	//		return error
	//	}
	//	go ScaleImage(source.Source)
	//	mlog.Info("json:%s",source.Source)
	//	return c.JSON(http.StatusOK,"{\"return\":\"ok\"}")
	//})
	//e.SetDebug(true)
	//go e.Run(standard.New(":4444"))
//}

//func ScaleImage(source string) {
//	cmd := exec.Command("cwebp", "-resize", "480", "320", string(source), "-o", "1.webp")
//	if err := cmd.Run(); err != nil {
//		mlog.Error(err)
//	}
//	mlog.Info("convert %s ok", source)
//}

//

package main
import (
	"github.com/kataras/iris"
	"github.com/jbrodriguez/mlog"
	"work/router"
	"github.com/golang/protobuf/proto"
	"work/myproto"
	"work/message"
)



func main() {
	mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)
	iris.Config.Websocket.Endpoint = "/ws"
	r := router.New()
	r.Map(1000,func(p proto.Message) []byte{
		createRoom := p.(*myproto.CreateRoomTos)
		mlog.Info("create room :%s",createRoom.RoomName)
		roomID := "1000"
		rData,_ := message.Marshal(&myproto.CreateRoomToc{RoomID:&roomID})
		return rData
	})
	iris.Websocket.OnConnection(func(c iris.WebsocketConnection){
		mlog.Info("Connect websocket")
		c.OnMessage(func(message []byte){
			mlog.Info("receive %v",message)
			sendMsg := r.DoRoute(&message)

			c.EmitMessage(sendMsg)
		})
	})

	iris.Listen(":8080")
}

//package main
//
//import (
//	"regexp"
//	"io/ioutil"
//	"os"
//
//)
//
//func main() {
//	// Compile the expression once, usually at init time.
//	// Use raw strings to avoid having to quote the backslashes.
//	//message[ \t]*[ \t\n]*(\/\/<(.*)>)*[ \t\n]*[ \t\n]*{[ \t\n]*[^}]*}
//	var validID = regexp.MustCompile("message[ \\t]+([\\w]+)[ \\t\\r\\n]*({\\/\\/<([\\w]+)>)*[ \\t\\n]*[ \\t\\n]*([ \\t\\r\\n]+[^}]*)}")
//	d,_ :=ioutil.ReadFile("./live/myproto.proto")
//	reg := validID.FindAllSubmatch(d,-1)
//
//	id2NameStr :="const id2Name = map[int64]interface{}{"
//	name2IDStr := "const name2IDStr = map[interface{}]int64{"
//	strTmp1 :="\n"
//	strTmp2 :=""
//	for _,match := range reg {
//
//		name := match[1]
//		id := match[3]
//		if strTmp1 != "\n" {
//			strTmp1 = strTmp1 +",\n"  + string(name) + "{}" + " : " + string(id)
//			strTmp2 = strTmp2 + ",\n"  + string(name) + "{}" + " : " + string(id)
//		} else {
//			strTmp1 = strTmp1 + string(id) + " : " +string(name) + "{}"
//			strTmp2 =  strTmp2 + string(name) + "{}" + " : " + string(id)
//		}
//	}
//		id2NameStr = id2NameStr + strTmp1 + "}\n"
//		name2IDStr = name2IDStr + strTmp2 + "}\n"
//
//	f, err := os.OpenFile("./live/live.pb.go", os.O_APPEND|os.O_WRONLY, 0600)
//	if err != nil {
//		panic(err)
//	}
//
//	defer f.Close()
//
//	if _, err = f.WriteString(id2NameStr + "\n" +name2IDStr); err != nil {
//		panic(err)
//	}
//
//}
//package main
//
//import (
//	"reflect"
//	"fmt"
//)
//
//func main() {
//	x := 6.4
//	p := reflect.ValueOf(&x)
//
//	fmt.Println("type of p:",p.Type())
//	fmt.Println("can set p:",p.CanSet())
//	v:=p.Elem()
//	fmt.Println("type of v:",v.Type())
//	fmt.Println("can set v:",v.CanSet())
//	v.SetFloat(7.1)
//	fmt.Println(v.Interface())
//	fmt.Println(v)
//
//	x1 := reflect.ValueOf(x)
//	fmt.Println("type of x:",x1.Type())
//	fmt.Println("can set x:",x1.CanSet())
//	v.SetFloat(791.1)
//	fmt.Println(x1.Interface())
//	fmt.Println(x1)
//}

//package main
//import (
//	"fmt"
//	"golang.org/x/net/http2"
//	"html"
//	"log"
//	"net/http"
//)
//func main() {
//	var srv http.Server
//	http2.VerboseLogs = true
//	srv.Addr = ":8080"
//	// This enables http2 support
//	http2.ConfigureServer(&srv, nil)
//	// Plain text test handler
//	// Open https://localhost:8080/randomtest
//	// in your Chrome Canary browser
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "Hi tester %q\n", html.EscapeString(r.URL.Path))
//		ShowRequestInfoHandler(w, r)
//	})
//	// Listen as https ssl server
//	// NOTE: WITHOUT SSL IT WONT WORK!!
//	// To self generate a test ssl cert/key you could go to
//	// http://www.selfsignedcertificate.com/
//	// or read the openssl manual
//	log.Fatal(srv.ListenAndServeTLS("localhost.cert","localhost.key"))
//}
//func ShowRequestInfoHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//	fmt.Fprintf(w, "Method: %s\n", r.Method)
//	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
//	fmt.Fprintf(w, "Host: %s\n", r.Host)
//	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
//	fmt.Fprintf(w, "RequestURI: %q\n", r.RequestURI)
//	fmt.Fprintf(w, "URL: %#v\n", r.URL)
//	fmt.Fprintf(w, "Body.ContentLength: %d (-1 means unknown)\n", r.ContentLength)
//	fmt.Fprintf(w, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
//	fmt.Fprintf(w, "TLS: %#v\n", r.TLS)
//	fmt.Fprintf(w, "\nHeaders:\n")
//	r.Header.Write(w)
//}