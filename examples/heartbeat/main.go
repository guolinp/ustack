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
		SetDataProcessor(ustack.NewHeartbeat().SetName("HB-Client").ForServer(false)).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		SetEventListener(func(event ustack.Event) {
			if event.Type == ustack.UStackEventHeartbeatLost {
				connection := event.Data.(ustack.TransportConnection)
				fmt.Println("connection:", connection.GetName(), "heartbeat lost")
			} else if event.Type == ustack.UStackEventHeartbeatRecover {
				connection := event.Data.(ustack.TransportConnection)
				fmt.Println("connection:", connection.GetName(), "heartbeat recover")
			}
		}).
		Run()
}

func server() {
	ustack.NewUStack().
		SetName("Server").
		SetDataProcessor(ustack.NewHeartbeat().SetName("HB-Server").ForServer(true)).
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
