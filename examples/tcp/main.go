package main

import (
	"fmt"
	"os"
	"time"

	"ustack"
)

var ustackName string = "UStack"
var dataServerAddress string = "127.0.0.1:1234"
var statServerAddress1 string = "127.0.0.1:8800"
var statServerAddress2 string = "127.0.0.1:8801"

func client() {
	ustack.NewUStack().
		SetName(ustackName).
		AddEndPoint(
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Millisecond * 100)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("1234567890"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("ACK:", epd.GetConnection().GetName(), string(epd.GetData().([]byte)))
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewStatCounter()).
		AppendDataProcessor(ustack.NewFrameDecoder().SetOption("CacheCapacity", 128)).
		AddTransport(
			ustack.NewTCPTransport("tcpDataClient").
				ForServer(false).
				SetAddress(dataServerAddress)).
		AddTransport(
			ustack.NewTCPTransport("tcpStatServer").
				ForServer(true).
				SetAddress(statServerAddress1)).
		Run()

}

func server() {
	ustack.NewUStack().
		SetName(ustackName).
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
						fmt.Println("REQ:", epd.GetConnection().GetName(), string(epd.GetData().([]byte)))

						endpoint.GetTxChannel() <- ustack.NewEndPointData(epd.GetConnection(), []byte("1:0987654321"))
						endpoint.GetTxChannel() <- ustack.NewEndPointData(epd.GetConnection(), []byte("2:0987654321"))
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewStatCounter()).
		AppendDataProcessor(ustack.NewFrameDecoder().SetOption("CacheCapacity", 128)).
		AddTransport(
			ustack.NewTCPTransport("tcpDataServer").
				ForServer(true).
				SetAddress(dataServerAddress)).
		AddTransport(
			ustack.NewTCPTransport("tcpStatServer").
				ForServer(true).
				SetAddress(statServerAddress2)).
		Run()
}

func update() {
	if len(os.Args) > 2 {
		ustackName = os.Args[2]
		fmt.Println("ustackName:", ustackName)
	}

	if len(os.Args) > 3 {
		dataServerAddress = os.Args[3]
		fmt.Println("dataServerAddress:", dataServerAddress)
	}

	if len(os.Args) > 4 {
		statServerAddress1 = os.Args[4]
		statServerAddress2 = os.Args[4]
		fmt.Println("statServerAddress:", statServerAddress1)
	}
}

func main() {
	if len(os.Args) > 1 {
		if fn, ok := map[string]func(){
			"-s": server,
			"-c": client,
		}[os.Args[1]]; ok {
			update()
			fn()
			time.Sleep(time.Second * 3600)
			return
		}
	}

	fmt.Println(os.Args[0], "<-s|-c|-h> [ustackName] [dataServerAddress] [statServerAddress]")
}
