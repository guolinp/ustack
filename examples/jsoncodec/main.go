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
		AddEndPoint(
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							user := &User{Name: "ZhangSan", Age: 40}
							fmt.Println("Send:", user.Name, user.Age)
									endpoint.GetTxChannel() <- ustack.NewEndPointData().
										SetConnection(connection).
										SetData(user)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AppendDataProcessor(ustack.NewJSONCodec(reflect.TypeOf(User{}))).
		AddTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func server() {
	ustack.NewUStack().
		SetName("Server").
		AddEndPoint(
			ustack.NewEndPoint("EP-Server", 0).
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
		AppendDataProcessor(ustack.NewJSONCodec(reflect.TypeOf(User{}))).
		AddTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func main() {
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

	fmt.Println(os.Args[0], "<-s|-c|-h>")
}
