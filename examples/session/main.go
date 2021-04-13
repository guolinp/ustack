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
			ustack.NewEndPoint("EP-Client1", 86).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Second * 2)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("EP1:86"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AddEndPoint(
			ustack.NewEndPoint("EP-Client2", 45).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Second * 3)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("EP2:45"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AppendDataProcessor(ustack.NewBytesCodec().SetName("PR-In-Client")).
		AppendDataProcessor(ustack.NewSessionResolver().SetName("PR-In-Client")).
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
			ustack.NewEndPoint("EP-Server1", 86).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("EP1: Receive: ", string(epd.GetData().([]byte)))
					})).
		AddEndPoint(
			ustack.NewEndPoint("EP-Server2", 45).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("EP2: Receive: ", string(epd.GetData().([]byte)))
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
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
