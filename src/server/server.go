package server

import (
	"ChatRoom/client"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type XServer struct {
	xListener   net.Listener
	port        string
	mux         sync.Mutex
	clients     map[int]*client.XClient
	connections int
	close       bool
}

func (XS *XServer) SetPort(p int) bool {
	if p < 4000 || p > 8000 {
		return false
	}
	XS.port = strconv.Itoa(p)
	return true
}
func (XS *XServer) portIsSet() bool {
	if XS.port == "" {
		fmt.Println("Please specify a port before starting the server!")
		return false
	}
	return true
}
func (XS *XServer) sockOpen() bool {
	xListener, err := net.Listen("tcp4", ":"+XS.port)
	XS.xListener = xListener
	if err != nil {
		fmt.Printf("Couldn't Open Socket\nReason:\n%s")
		fmt.Println(err)
		return false
	}
	return true
}
func (XS *XServer) serveClient(id int) {
	fmt.Printf("%s Connected!\n", XS.clients[id].Client.RemoteAddr().String())
	for !XS.clients[id].Inactive {
		cmd, err := bufio.NewReader(XS.clients[id].Client).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := strings.TrimSpace(string(cmd))
		fmt.Printf("%s: %s\n", XS.clients[id].Client.RemoteAddr().String(), msg)
		msgarr := strings.Fields(msg)
		if len(msgarr) > 0 {
			switch msgarr[0] {
			case "/Disconnect":
				XS.clients[id].Client.Write([]byte(string("Good bye!\n")))
				XS.clients[id].Inactive = true
			case "/Nick":
				XS.clients[id].Nickname = msgarr[1]
			case "/Say":
				var tmpstr string
				for i, str := range msgarr {
					if i != 0 {
						tmpstr += str + " "
					}
				}
				for _, user := range XS.clients {
					user.Client.Write([]byte(string(XS.clients[id].Nickname + ": " + tmpstr + "\n")))
				}
			default:
				XS.clients[id].Client.Write([]byte(string("Unrecognized Command!\n")))
			}
		}
	}
	XS.clients[id].Client.Close()
	XS.mux.Lock()
	delete(XS.clients, id)
	defer XS.mux.Unlock()
}
func (XS *XServer) newClient(c net.Conn) {
	XS.mux.Lock()
	XS.connections++
	XS.clients[XS.connections] = &client.XClient{c, "", false}
	go XS.serveClient(XS.connections)
	defer XS.mux.Unlock()
}
func (XS *XServer) listen() {
	for !XS.close {
		c, err := XS.xListener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go XS.newClient(c)
	}
}
func (XS *XServer) Start() bool {
	XS.clients = make(map[int]*client.XClient)
	fmt.Println("Starting Server...")
	if !XS.portIsSet() || !XS.sockOpen() {
		return false
	}
	fmt.Println("Server successfully started!")
	XS.listen()
	defer XS.xListener.Close()
	return true
}
