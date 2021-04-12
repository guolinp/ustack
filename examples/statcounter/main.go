package main

import (
	"fmt"
	"os"
	"time"

	"ustack"
)

func client() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Millisecond * 500)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("1234"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewBytesCodec()).
		SetDataProcessor(ustack.NewStatCounter()).
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
						fmt.Println("receive:", string(epd.GetData().([]byte)))
					})).
		SetDataProcessor(ustack.NewBytesCodec()).
		SetDataProcessor(ustack.NewStatCounter()).
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
