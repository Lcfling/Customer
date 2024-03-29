package udpsok

import (
	"fmt"
	log "github.com/golang/glog"
	"net"
	"time"
)

type Client struct {
	Connection
	f func(e bool)
}

func NewClient(conn *net.UDPConn, sn int64, addr *net.UDPAddr) *Client {
	client := new(Client)
	client.wt = make(chan *Message, 300)
	client.rt = make(chan []byte, 300)
	client.addr = addr
	client.sn = sn
	client.conn = conn
	client.lastime = time.Now().Unix()
	client.Run()
	return client
}
func (client *Client) UpdateAddr(addr *net.UDPAddr) {
	client.addr = addr
	client.lastime = time.Now().Unix()
}
func (client *Client) RemoveClient() {
	mutex.Lock()
	defer mutex.Unlock()
	close(client.wt)
	close(client.rt)
	delete(Clients, client.sn)
}
func (client *Client) Write() {
	running := true
	for running {
		select {
		case msg := <-client.wt:
			if msg == nil {
				//client.close()
				running = false
				log.Infof("client:%d socket closed", client.sn)
				return
			}
			client.SendMessage(msg)
		}
	}
}
func (client *Client) Read() {
	running := true
	for running {
		select {
		case data := <-client.rt:
			if data == nil {
				fmt.Println("client read rt nil exit")
				return
			}
			msg := client.read(data)

			client.HandleMessage(msg)
		}
	}
}
func (client *Client) Check() {
	for true {
		tiker := time.NewTicker(time.Second * 1)
		<-tiker.C
		log.Info("client check  running")
		if time.Now().Unix()-client.lastime > 10 {
			client.RemoveClient()
			return
		}
	}

}

func (client *Client) Run() {
	//go client.Check()
	go client.Write()
	go client.Read()
}
func (client *Client) HandleMessage(msg *Message) {
	log.Info("msg cmd:", msg.cmd)
	switch msg.cmd {

	case MSG_OPEN:
		client.HandleOpen(msg.body.(*OpenMessage))
	case MSG_PING:
		client.HandlePing()

	}
}

func (client *Client) HandlePing() {
	client.lastime = time.Now().Unix()
}
func (client *Client) HandleOpen(open *OpenMessage) {
	fmt.Println("分发消2息")
	if open.index == 1 {
		client.f(true)
	} else {
		client.f(false)
	}
}
