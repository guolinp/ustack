package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"

	"ustack"
)

type Shower interface {
	Show() string
}

// Person ...
type Person struct {
	Name string
	Age  int
}

func (p *Person) Show() string {
	return fmt.Sprintf("Person: Name: %s Age: %d", p.Name, p.Age)
}

// Dog ...
type Dog struct {
	Name string
}

func (d *Dog) Show() string {
	return fmt.Sprintf("Dog: Name: %s", d.Name)
}

func doInit() {
	gob.Register(&Person{})
	gob.Register(&Dog{})
}

func encoder(m interface{}, w io.Writer) error {
	fmt.Println("Call encoder")

	return gob.NewEncoder(w).Encode(m)
}

func decoder(r io.Reader) (interface{}, error) {
	fmt.Println("call decoder")

	var s Shower
	err := gob.NewDecoder(r).Decode(&s)
	return s, err
}

func client() {
	ustack.NewUStack().
		SetName("Client").
		AddEndPoint(
			ustack.NewEndPoint("EP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)

							var p Shower = &Person{Name: "LiSi", Age: 10}
							fmt.Println("Send:", p.Show())
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, &p)

							var d Shower = &Dog{Name: "Tom"}
							fmt.Println("Send:", d.Show())
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, &d)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		AppendDataProcessor(ustack.NewGenericCodec(doInit, encoder, decoder)).
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
			ustack.NewEndPoint("EP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						shower := epd.GetData().(Shower)
						fmt.Println("Receive:", shower.Show())
					})).
		AppendDataProcessor(ustack.NewGenericCodec(doInit, encoder, decoder)).
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
