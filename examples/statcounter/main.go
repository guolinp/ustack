package main

import (
	"fmt"
	"os"
	"time"

	"ustack"
)

var statServerAddress1 string = "127.0.0.1:8800"
var statServerAddress2 string = "127.0.0.1:8801"

func client() {
	ustack.NewUStack().
		SetName("tcpStatClient").
		AppendDataProcessor(ustack.NewStatCounter().SetOption("Collect.IntervalInSecond", 2)).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(
			ustack.NewTCPTransport("tcpToStatServer1").
				ForServer(false).
				SetAddress(statServerAddress1)).
		AddTransport(
			ustack.NewTCPTransport("tcpToStatServer2").
				ForServer(false).
				SetAddress(statServerAddress2)).
		Run()

}

func update() {
	if len(os.Args) > 2 {
		statServerAddress1 = os.Args[2]
		fmt.Println("statServerAddress1:", statServerAddress1)
	}

	if len(os.Args) > 3 {
		statServerAddress2 = os.Args[3]
		fmt.Println("statServerAddress2:", statServerAddress1)
	}
}

func main() {
	if len(os.Args) > 1 {
		if fn, ok := map[string]func(){
			"-c": client,
		}[os.Args[1]]; ok {
			update()
			fn()
			time.Sleep(time.Second * 3600)
			return
		}
	}

	fmt.Println(os.Args[0], "<-s|-c|-h> [statServerAddress1] [statServerAddress2]")
}
