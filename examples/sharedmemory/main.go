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
									data := "SharedMemory 1234"
									fmt.Println("Send:", data)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte(data))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AddTransport(
			ustack.NewSharedMemoryTransport("SharedMemoryClient").
				ForServer(false).
				SetAddress("SharedMemory")).
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
						fmt.Println("Receive:", string(epd.GetData().([]byte)))
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AddTransport(
			ustack.NewSharedMemoryTransport("SharedMemoryServer").
				ForServer(true).
				SetAddress("SharedMemory")).
		Run()
}

func main() {
	server()
	client()

	time.Sleep(time.Second * 3600)
}
