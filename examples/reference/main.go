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
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Millisecond * 100)
									data := "1234567890"
									fmt.Println("Send:", data)
									endpoint.GetTxChannel() <- ustack.NewEndPointData().
										SetConnection(connection).
										SetData([]byte(data))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AddTransport(
			ustack.NewReferenceTransport("Client").
				ForServer(false).
				SetAddress("UnifyString")).
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
						fmt.Println("Receive:", string(epd.GetData().([]byte)))
					})).
		AddTransport(
			ustack.NewReferenceTransport("Server").
				ForServer(true).
				SetAddress("UnifyString")).
		Run()
}

func main() {
	server()
	client()

	time.Sleep(time.Second * 3600)
}
