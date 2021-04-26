package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"ustack"
)

var ustackName string = "UStack"
var dataServerAddress string = "127.0.0.1:1234"
var packageSize string = "64"
var packageData []byte

var TxCount uint64 = 0
var RxCount uint64 = 0

func ShowSpeed() {
	var seconds uint64 = 0

	for {
		time.Sleep(2 * time.Second)
		seconds += 2

		fmt.Printf("Second: %-4d Count: Tx %-12d Rx %-12d  Speed: Tx %-12d Rx %-12d\n",
			seconds, TxCount, RxCount, TxCount/seconds, RxCount/seconds)
	}
}

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
									//time.Sleep(time.Millisecond * 1)
									endpoint.GetTxChannel() <- ustack.NewEndPointData().
										SetConnection(connection).
										SetData(packageData)
									TxCount++
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						RxCount++
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(
			ustack.NewTCPTransport("tcpDataClient").
				ForServer(false).
				SetAddress(dataServerAddress)).
		Run()

	ShowSpeed()
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
						//fmt.Println("ACK:", epd.GetConnection().GetName(), string(epd.GetData().([]byte)))
						RxCount++
						endpoint.GetTxChannel() <- ustack.NewEndPointData().
							SetConnection(epd.GetConnection()).
							SetData(packageData)
						TxCount++
					})).
		AppendDataProcessor(ustack.NewBytesCodec()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(
			ustack.NewTCPTransport("tcpDataServer").
				ForServer(true).
				SetAddress(dataServerAddress)).
		Run()

	ShowSpeed()
}

func update() {
	if len(os.Args) > 2 {
		dataServerAddress = os.Args[2]
		fmt.Println("dataServerAddress:", dataServerAddress)
	}

	if len(os.Args) > 3 {
		packageSize = os.Args[3]
		fmt.Println("packageSize:", packageSize)
	}
}

func makePackageData() {
	size, err := strconv.Atoi(packageSize)
	if err != nil {
		size = 64
	}

	packageData = make([]byte, size)
	for i := 0; i < size; i++ {
		packageData[i] = 0x41
	}
	fmt.Println("package data size:", size)
}

func main() {
	if len(os.Args) > 1 {
		if fn, ok := map[string]func(){
			"-s": server,
			"-c": client,
		}[os.Args[1]]; ok {
			update()
			makePackageData()
			fn()
			ShowSpeed()
			time.Sleep(time.Second * 3600)
			return
		}
	}

	fmt.Println(os.Args[0], "<-s|-c|-h> [dataServerAddress] [packageSize]")
}
