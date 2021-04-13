package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"ustack"
)

// User ...
type User struct {
	Name string
	Age  int
}

func encoder(message interface{}, ub *ustack.UBuf) error {
	fmt.Println("Call encoder")

	return gob.NewEncoder(ub).Encode(message)
}

func decoder(ub *ustack.UBuf) (interface{}, error) {
	fmt.Println("call decoder")

	var user User = User{}
	err := gob.NewDecoder(ub).Decode(&user)
	return user, err
}

func client() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							user := &User{Name: "LiSi", Age: 10}
							fmt.Println("Send:", user.Name, user.Age)
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, user)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AppendDataProcessor(ustack.NewGenericCodec(encoder, decoder)).
		AddTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()

}

func server() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("EP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						user := epd.GetData().(User)
						fmt.Println("Receive:", user.Name, user.Age)
					})).
		AppendDataProcessor(ustack.NewGenericCodec(encoder, decoder)).
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
