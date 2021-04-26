package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"ustack"
)

var centerServiceAddress string = "127.0.0.1:1234"
var helloServiceAddress string = "127.0.0.1:5678"

/**********************************************************************************************
 * Service center
 **********************************************************************************************/
func startCenterServer() {
	services := map[string]string{
		"findHelloAddress": helloServiceAddress,
	}

	ustack.NewUStack().
		SetName("CenterServer").
		AddEndPoint(
			ustack.NewEndPoint("CenterServer:EP-0", 0).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						request := epd.GetData().(string)

						fmt.Println("CenterServer: Receive:", request)

						if strings.HasSuffix(request, "reg") {
							info := strings.Split(request, ":")
							services[info[0]] = info[1]
							fmt.Println("CenterServer: reg:", info[0], services[info[0]])
						} else {
							addr, ok := services[request]
							if ok {
								endpoint.GetTxChannel() <- ustack.NewEndPointData().
									SetConnection(epd.GetConnection()).
									SetData(addr)
								fmt.Println("CenterServer: send:", addr)
							} else {
								fmt.Println("CenterServer: not found service", request, "address")
							}
						}
					})).
		AppendDataProcessor(ustack.NewStringCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(
			ustack.NewTCPTransport("tcpCenterServer").
				ForServer(true).
				SetAddress(centerServiceAddress)).
		Run()
}

/**********************************************************************************************
 * Service provider
 **********************************************************************************************/
func startHelloServer() {
	ustack.NewUStack().
		SetName("HelloServer").
		AddEndPoint(
			ustack.NewEndPoint("HelloServer:EP-0", 0).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						request := epd.GetData().(string)

						fmt.Println("HelloServer: Receive:", request)

						if request == "hello" {
							endpoint.GetTxChannel() <- ustack.NewEndPointData().
								SetConnection(epd.GetConnection()).
								SetData("world")
							fmt.Println("HelloServer: Send: world")
						}
					})).
		AppendDataProcessor(ustack.NewStringCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(
			ustack.NewTCPTransport("tcpHelloServer").
				ForServer(true).
				SetAddress(helloServiceAddress)).
		Run()
}

/**********************************************************************************************
 * App Session
 **********************************************************************************************/
type AppSession struct {
	u     ustack.UStack
	ep    ustack.EndPoint
	c     ustack.TransportConnection
	ready chan int
}

func NewAppSession(address string, forServer bool) *AppSession {
	session := &AppSession{
		c:     nil,
		ready: make(chan int, 1),
	}

	session.ep = ustack.NewEndPoint("EP0", 0).
		SetEventListener(
			func(endpoint ustack.EndPoint, event ustack.Event) {
				if event.Type == ustack.UStackEventNewConnection {
					session.c = event.Data.(ustack.TransportConnection)
					session.ready <- 1
				}
			})

	session.u = ustack.NewUStack().
		AddEndPoint(session.ep).
		AppendDataProcessor(ustack.NewStringCodec()).
		AppendDataProcessor(ustack.NewSessionResolver()).
		AppendDataProcessor(ustack.NewFrameDecoder()).
		AddTransport(
			ustack.NewTCPTransport("TP").
				ForServer(forServer).
				SetAddress(address))

	return session
}

func (s *AppSession) WaitReady() {
	<-s.ready
}

func (s *AppSession) Read() string {
	epd := <-s.ep.GetRxChannel()
	s.c = epd.GetConnection()
	return epd.GetData().(string)
}

func (s *AppSession) Write(data string) {
	s.ep.GetTxChannel() <- ustack.NewEndPointData().
		SetConnection(s.c).
		SetData(data)
}

func (s *AppSession) Run(fn func(s *AppSession)) {
	s.u.Run()
	s.WaitReady()
	go fn(s)
}

/**********************************************************************************************
 * App Example
 **********************************************************************************************/

func startApp() {
	// Service discovery session
	NewAppSession(centerServiceAddress, false).Run(func(s *AppSession) {
		s.Write("findHelloAddress")
		addr := s.Read()
		fmt.Println("addr:", addr)

		// create new session for 'hello' services
		NewAppSession(addr, false).Run(func(s *AppSession) {
			for {
				s.Write("hello")
				fmt.Println(s.Read())
				time.Sleep(time.Second)
			}
		})
	})
}

func main() {
	if len(os.Args) > 1 {
		if fn, ok := map[string]func(){
			"-s1": startCenterServer,
			"-s2": startHelloServer,
			"-c":  startApp,
		}[os.Args[1]]; ok {
			fn()
			time.Sleep(time.Second * 3600)
			return
		}
	}

	fmt.Println(os.Args[0], "<-s1|-s2|-c|-h>")
}
