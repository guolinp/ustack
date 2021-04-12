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
		SetEventListener(func(event ustack.Event) {
			if event.Type == ustack.UStackEventHeartbeatLost {
				connection := event.Data.(ustack.TransportConnection)
				fmt.Println("connection:", connection.GetName(), "heartbeat lost")
			} else if event.Type == ustack.UStackEventHeartbeatRecover {
				connection := event.Data.(ustack.TransportConnection)
				fmt.Println("connection:", connection.GetName(), "heartbeat recover")
			}
		}).
		AppendDataProcessor(ustack.NewHeartbeat().ForServer(false)).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func server() {
	ustack.NewUStack().
		SetName("Server").
		AppendDataProcessor(ustack.NewHeartbeat().ForServer(true)).
		SetTransport(
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
