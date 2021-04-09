package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"ustack"
)

func client1() {
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

func server1() {
	ustack.NewUStack().
		SetName("Server").
		SetDataProcessor(ustack.NewHeartbeat().SetName("HB-Server").ForServer(true)).
		SetTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func client2() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Millisecond * 1000)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("1234"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewBasicCodec().SetName("PR-In-Client")).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()

}

func server2() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("receive:", string(epd.GetData().([]byte)))
					})).
		SetDataProcessor(ustack.NewBasicCodec().SetName("PR-In-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func client3() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client1", 86).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Second * 3)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("EP1:86"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client2", 45).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Second * 7)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("EP2:45"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewBasicCodec().SetName("PR-In-Client")).
		SetDataProcessor(ustack.NewSessionResolver().SetName("PR-In-Client")).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func server3() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server1", 86).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("EP1: Recv: ", string(epd.GetData().([]byte)))
					})).
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server2", 45).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("EP2: Recv: ", string(epd.GetData().([]byte)))
					})).
		SetDataProcessor(ustack.NewBasicCodec().SetName("PR-In-Server")).
		SetDataProcessor(ustack.NewSessionResolver().SetName("PR-In-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

// User ...
type User struct {
	Name string
	Age  int
}

func client4() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							user := &User{Name: "ZhangSan", Age: 40}
							fmt.Println("Client Send:", user.Name, user.Age)
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, user)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewJSONCodec(reflect.TypeOf(User{})).SetName("JSONCODEC-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()

}

func server4() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						user := epd.GetData().(*User)
						fmt.Println("Server receive:", user.Name, user.Age)
					})).
		SetDataProcessor(ustack.NewJSONCodec(reflect.TypeOf(User{})).SetName("JSONCODEC-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func client5() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							user := &User{Name: "LiSi", Age: 10}
							fmt.Println("Client Send:", user.Name, user.Age)
							endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, user)
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewGOBCodec(reflect.TypeOf(User{})).SetName("GOBCODEC-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpClient").
				ForServer(false).
				SetAddress("127.0.0.1:1234")).
		Run()

}

func server5() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						user := epd.GetData().(*User)
						fmt.Println("Server receive:", user.Name, user.Age)
					})).
		SetDataProcessor(ustack.NewGOBCodec(reflect.TypeOf(User{})).SetName("GOBCODEC-Server")).
		SetTransport(
			ustack.NewTCPTransport("tcpServer").
				ForServer(true).
				SetAddress("127.0.0.1:1234")).
		Run()
}

func client6() {
	ustack.NewUStack().
		SetName("Client").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Client", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventNewConnection {
							connection := event.Data.(ustack.TransportConnection)
							go func() {
								for {
									time.Sleep(time.Millisecond * 1000)
									endpoint.GetTxChannel() <- ustack.NewEndPointData(connection, []byte("1234"))
								}
							}()
						} else if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					})).
		SetDataProcessor(ustack.NewBasicCodec().SetName("PR-In-Client")).
		SetTransport(
			ustack.NewUDSTransport("udsClient").
				ForServer(false).
				SetAddress("/tmp/uds.socket")).
		Run()

}

func server6() {
	ustack.NewUStack().
		SetName("Server").
		SetEndPoint(
			ustack.NewEndPoint("AppEP-Server", 0).
				SetEventListener(
					func(endpoint ustack.EndPoint, event ustack.Event) {
						if event.Type == ustack.UStackEventConnectionClosed {
							os.Exit(1)
						}
					}).
				SetDataListener(
					func(endpoint ustack.EndPoint, epd ustack.EndPointData) {
						fmt.Println("receive:", string(epd.GetData().([]byte)))
					})).
		SetDataProcessor(ustack.NewBasicCodec().SetName("PR-In-Server")).
		SetTransport(
			ustack.NewUDSTransport("udsServer").
				ForServer(true).
				SetAddress("/tmp/uds.socket")).
		Run()
}

func main() {
	map[string]func(){
		"s1": server1,
		"c1": client1,
		"s2": server2,
		"c2": client2,
		"s3": server3,
		"c3": client3,
		"s4": server4,
		"c4": client4,
		"s5": server5,
		"c5": client5,
		"s6": server6,
		"c6": client6,
	}[os.Args[1]]()

	time.Sleep(time.Second * 3600)
}
