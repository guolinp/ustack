package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"ustack"
)

// User ...
type User struct {
	Name string
	Age  int
}

func client() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							user := &User{Name: "ZhangSan", Age: 40}
							fmt.Println("Client Send:", user.Name, user.Age)
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, user)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewJSONCodec(reflect.TypeOf(User{})).SetName("JSONCODEC-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func server() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						user := epd.GetData().(*User)
						fmt.Println("Server receive:", user.Name, user.Age)
					})).
		SetDataProcessor(ustack.NewJSONCodec(reflect.TypeOf(User{})).SetName("JSONCODEC-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func main() {
	help := func() { fmt.Println(os.Args[0], "<-s|-c|-h>") }
	if len(os.Args) > 1 {
		if fn, ok := map[string]func(){
			"-s": server,
			"-c": client,
		}[os.Args[1]]; ok {
			fn()
			time.Sleep(time.Second * 3600)
			return
		}
	}

	help()
}
