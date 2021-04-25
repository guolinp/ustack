package main

import (
	"fmt"
	"os"
	"time"

	"ustack"
)

var dataServerAddress string = "127.0.0.1:1234"

var tcpForClient ustack.Transport = ustack.NewTCPTransport("tcpClient").
	ForServer(false).
	SetAddress(dataServerAddress)

var tcpForServer ustack.Transport = ustack.NewTCPTransport("tcpServer").
	ForServer(true).
	SetAddress(dataServerAddress)

var udsForClient ustack.Transport = ustack.NewUDSTransport("udsClient").
	ForServer(false).
	SetAddress("/tmp/gouds.socket")

var udsForServer ustack.Transport = ustack.NewUDSTransport("udsServer").
	ForServer(true).
	SetAddress("/tmp/gouds.socket")

var epForClient ustack.EndPoint = ustack.NewEndPoint("NewEP-Client", 10).
	SetEventListener(
		func(endpoint ustack.EndPoint, event ustack.Event) {
			if event.Type == ustack.UStackEventNewConnection {
				connection := event.Data.(ustack.TransportConnection)
				fmt.Println("UStackEventNewConnection", connection.GetName())

				go func() {
					for {
						if connection.Closed() {
							return
						}
						time.Sleep(time.Millisecond * 1000)
						data := []byte("abcd")
						fmt.Println("EP10-Send:", string(data))
						endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, data)
					}
				}()
			}
		})

var epForServer ustack.EndPoint = ustack.NewEndPoint("NewEP-Server", 10).
	SetDataListener(
		func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
			fmt.Println("EP10-Receive:", string(epd.GetData().([]byte)))
		})

var u ustack.UStack
var conn int = 0

func client() {
	u = ustack.NewUStack().
		SetName("Client").
		AddEndPoint(
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							fmt.Println("UStackEventNewConnection", connection.GetName())

							conn++

							if conn == 1 {
								fmt.Println("add new transport", tcpForClient.GetName())
								u.AddEndPoint(epForClient)
								u.AddTransport(tcpForClient)
							}

							go func(c int) {
								var data []byte
								if c == 1 {
									data = []byte("uds:1234")
								} else {
									data = []byte("tcp:5678")
								}
								for {
									if connection.Closed() {
										return
									}
									time.Sleep(time.Millisecond * 1000)
									fmt.Println("EP0-Send:", string(data))
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, data)
								}
							}(conn)

							if conn == 2 {
								go func(c int) {
									time.Sleep(time.Second * 5)
									fmt.Println("DeleteTransport")
									u.DeleteTransport(tcpForClient)
									time.Sleep(time.Second * 5)
									fmt.Println("DeleteEndPoint")
									u.DeleteEndPoint(epForClient)
								}(conn)
							}
						} else if event.Type == ustack.UStackEventConnectionClosed {
							connection := event.Data.(ustack.TransportConnection)
							fmt.Println("UStackEventConnectionClosed", connection.GetName())
							conn--
							if conn == 0 {
								os.Exit(1)
							}
						}
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(udsForClient).
		Run()

}

func server() {
	u = ustack.NewUStack().
		SetName("Server").
		AddEndPoint(
			ustack.NewEndPoint("EP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							fmt.Println("UStackEventNewConnection", connection.GetName())

							conn++

							if conn == 1 {
								fmt.Println("add new transport", tcpForServer.GetName())
								u.AddTransport(tcpForServer)
								u.AddEndPoint(epForServer)
							}
						}

						if event.Type == ustack.UStackEventConnectionClosed {
							connection := event.Data.(ustack.TransportConnection)
							fmt.Println("UStackEventConnectionClosed", connection.GetName())
							conn--
							if conn == 0 {
								os.Exit(1)
							}
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("EP0-Receive:", string(epd.GetData().([]byte)))
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(udsForServer).
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
