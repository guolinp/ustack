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
	err := gob.NewEncoder(ub).Encode(message)
	fmt.Println("call encoder")
	return err
}

func decoder(ub *ustack.UBuf) (interface{}, error) {
	var user User = User{}
	err := gob.NewDecoder(ub).Decode(&user)
	fmt.Println("call decoder")
	if err != nil {
		return nil, err
	}
	return user, nil
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
							fmt.Println("Client Send:", user.Name, user.Age)
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, user)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewGenericCodec(encoder, decoder)).
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
						fmt.Println("Server receive:", user.Name, user.Age)
					})).
		SetDataProcessor(ustack.NewGenericCodec(encoder, decoder)).
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
