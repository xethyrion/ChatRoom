package client

import "net"

type XClient struct {
	Client        net.Conn
	Nickname      string
	Inactive      bool
	Authenticated bool
}
