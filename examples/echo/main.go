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
		AddEndPoint(
			ustack.NewEndPoint("EP1", 1).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								data := []byte("1234")
								for {
									time.Sleep(time.Millisecond * 1000)
									fmt.Println(endpoint.GetName(), "Send:", string(data))

									endpoint.GetTxChannel() <- ustack.NewEndPointData().
										SetConnection(connection).
										SetDestinationSession(2).
										SetData(data)
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AddEndPoint(
			ustack.NewEndPoint("EP2", 2).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println(endpoint.GetName(), "Receive:", string(epd.GetData().([]byte)))
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
		AddTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()

}

func server() {
	ustack.NewUStack().
		SetName("Server").
		AppendDataProcessor(ustack.NewEcho()).
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
