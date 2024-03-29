/**
 * Copyright (c) 2014-2015, GoBelieve
 * All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

package socket

import (
	"github.com/Lcfling/Customer/models/users"
	"math"
	"net"
	"strconv"
)
import "time"
import "sync/atomic"
import log "github.com/golang/glog"
import "container/list"
import "crypto/tls"
import "fmt"

type Client struct {
	Connection //必须放在结构体首部
	*PeerClient
	*RoomClient
	Handler   map[string]func(c *Client, r *Request)
	public_ip int32
	Hts       HttpSocket
}

func NewClient(conn interface{}) *Client {

	client := new(Client)
	//初始化Connection
	client.conn = conn // conn is net.Conn or engineio.Conn
	if net_conn, ok := conn.(net.Conn); ok {
		addr := net_conn.LocalAddr()
		if taddr, ok := addr.(*net.TCPAddr); ok {
			ip4 := taddr.IP.To4()
			client.public_ip = int32(ip4[0])<<24 | int32(ip4[1])<<16 | int32(ip4[2])<<8 | int32(ip4[3])
		}
	}
	client.wt = make(chan *Message, 300)
	client.lwt = make(chan int, 1) //only need 1
	//'10'对于用户拥有非常多的超级群，读线程还是有可能会阻塞
	client.pwt = make(chan []*Message, 10)
	client.messages = list.New()
	//atomic.AddInt64(&server_summary.nconnections, 1)
	client.PeerClient = &PeerClient{&client.Connection}
	client.RoomClient = &RoomClient{Connection: &client.Connection}
	client.Handler = make(map[string]func(c *Client, r *Request))
	HandleSocket(client)
	return client
}

func handle_client(conn net.Conn) {
	log.Infoln("handle new connection, remote address:", conn.RemoteAddr())
	client := NewClient(conn)
	client.Run()
}

func handle_ssl_client(conn net.Conn) {
	log.Infoln("handle new ssl connection,  remote address:", conn.RemoteAddr())
	client := NewClient(conn)
	client.Run()
}

func Listen(f func(net.Conn), port int) {
	listen_addr := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen("tcp", listen_addr)
	if err != nil {
		log.Errorf("listen err:%s", err)
		return
	}
	tcp_listener, ok := listen.(*net.TCPListener)
	if !ok {
		log.Error("listen err")
		return
	}

	for {
		client, err := tcp_listener.AcceptTCP()
		if err != nil {
			log.Errorf("accept err:%s", err)
			return
		}
		f(client)
	}
}

func ListenClient() {
	Listen(handle_client, config.port)
}

func ListenSSL(port int, cert_file, key_file string) {
	cert, err := tls.LoadX509KeyPair(cert_file, key_file)
	if err != nil {
		log.Fatal("load cert err:", err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	addr := fmt.Sprintf(":%d", port)
	listen, err := tls.Listen("tcp", addr, config)
	if err != nil {
		log.Fatal("ssl listen err:", err)
	}

	log.Infof("ssl listen...")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("ssl accept err:", err)
		}
		handle_ssl_client(conn)
	}
}

func (client *Client) Read() {
	for {
		tc := atomic.LoadInt32(&client.tc)
		if tc > 0 {
			log.Infof("quit read goroutine, client:%d write goroutine blocked", client.uid)
			client.HandleClientClosed()
			break
		}

		t1 := time.Now().Unix()
		msg := client.read()
		t2 := time.Now().Unix()
		if t2-t1 > 6*60 {
			log.Infof("client:%d socket read timeout:%d %d", client.uid, t1, t2)
		}
		if msg == nil {
			log.Infof("msg = nil  HandleClientClosed ------")
			client.HandleClientClosed()
			break
		}
		log.Infof("Read message ------")
		log.Info(msg)
		client.HandleMessage(msg)
		t3 := time.Now().Unix()
		if t3-t2 > 2 {
			log.Infof("client:%d handle message is too slow:%d %d", client.uid, t2, t3)
		}
	}
}

func (client *Client) RemoveClient() {
	route := app_route.FindRoute(client.appid)
	if route == nil {
		log.Warning("can't find app route")
		return
	}
	route.RemoveClient(client)
	if client.room_id > 0 {
		route.RemoveRoomClient(client.room_id, client)
	}
}

func (client *Client) HandleClientClosed() {
	//atomic.AddInt64(&server_summary.nconnections, -1)
	if client.uid > 0 {
		//atomic.AddInt64(&server_summary.nclients, -1)
		atomic.StoreInt32(&client.closed, 1)
		client.RemoveClient()
		//quit when write goroutine received
		client.wt <- nil

		client.RoomClient.Logout()
		client.PeerClient.Logout()
		//go ChangeUserOnline(client.uid, 0)
		SetOnline(client.uid, 0)
	}
}

func (client *Client) HandleMessage(msg *Message) {
	log.Info("msg cmd:", Command(msg.cmd))
	switch msg.cmd {
	case MSG_AUTH_TOKEN:
		client.HandleAuthToken(msg.body.(*AuthenticationToken), msg.version)
	case MSG_ACK:
		client.HandleACK(msg.body.(*MessageACK))
	case MSG_PING:
		client.HandlePing()
	case MSG_HTTP_MSG:
		client.HandleHttp(msg)
	}

	client.PeerClient.HandleMessage(msg)
	client.RoomClient.HandleMessage(msg)
}
func (client *Client) HttpWrite(body string) {
	client.Hts.body = body
	msg := &HttpSocket{}
	msg.url = client.Hts.url
	msg.header = client.Hts.header
	msg.body = client.Hts.body
	vmsg := &Message{cmd: MSG_HTTP_MSG, seq: 0, version: client.version, flag: 0, body: msg}
	log.Infof("client:%d Write message", client.uid)
	log.Info(vmsg)
	client.send(vmsg)
}

func (client *Client) HandleHttp(message *Message) {
	msg := message.body.(*HttpSocket)
	if _, exist := client.Handler[msg.url]; !exist {
		//请求的方法不存在 请再次将不存在方法的请求地址给记录下来

		client.close()
		return
	}
	request := new(Request)
	request.url = msg.url
	client.Hts.url = msg.url
	//request.ip=client.public_ip
	request.header = msg.header
	client.Hts.header = msg.header
	request.body = msg.body

	client.Handler[msg.url](client, request)
	//client.EnqueueMessage(m)

}

func (client *Client) AuthToken(uid int64, token string) (int64, int64, int, bool, error) {
	appid, utoken, err := LoadUserAccessToken(uid)

	if utoken != token || utoken == "" {
		return 0, 0, 0, false, err
	}
	//appid:=1
	//uid:=2
	if err != nil {
		return 0, 0, 0, false, err
	}

	/*forbidden, notification_on, err := GetUserPreferences(appid, uid)
	if err != nil {
		return 0, 0, 0, false, err
	}*/

	return appid, uid, 0, true, nil
}

func (client *Client) HandleAuthToken(login *AuthenticationToken, version int) {
	if client.uid > 0 {
		log.Info("repeat login")
		return
	}

	var err error
	div, _ := strconv.ParseInt(login.device_id, 10, 64)
	appid, uid, fb, on, err := client.AuthToken(div, login.token)
	if err != nil {
		log.Infof("auth token:%s err:%s", login.token, err)
		msg := &Message{cmd: MSG_AUTH_STATUS, version: version, body: &AuthenticationStatus{1}}
		client.EnqueueMessage(msg)
		return
	}
	if uid == 0 {
		log.Info("auth token uid==0")
		msg := &Message{cmd: MSG_AUTH_STATUS, version: version, body: &AuthenticationStatus{1}}
		client.EnqueueMessage(msg)
		return
	}
	/*if login.platform_id != PLATFORM_WEB && len(login.device_id) > 0 {
		client.device_ID, err = GetDeviceID(login.device_id, int(login.platform_id))
		if err != nil {
			log.Info("auth token uid==0")
			msg := &Message{cmd: MSG_AUTH_STATUS, version:version, body: &AuthenticationStatus{1}}
			client.EnqueueMessage(msg)
			return
		}
	}*/

	//is_mobile := login.platform_id == PLATFORM_IOS || login.platform_id == PLATFORM_ANDROID
	/*online := true
	if on && !is_mobile {
		online = false
	}*/
	//吧原来的连接断开
	route := app_route.FindRoute(appid)
	if route != nil {
		clients := route.FindClientSet(uid)
		if len(clients) != 0 {
			for c, _ := range clients {
				c.close()
			}
		}
	}

	client.appid = appid
	client.uid = uid
	client.forbidden = int32(fb)
	client.notification_on = on
	client.online = true
	client.version = version
	client.device_id = login.device_id
	client.platform_id = login.platform_id
	client.tm = time.Now()
	log.Infof("auth token:%s appid:%d uid:%d device id:%s:%d forbidden:%d notification on:%t online:%t",
		login.token, client.appid, client.uid, client.device_id,
		client.device_ID, client.forbidden, client.notification_on, client.online)

	msg := &Message{cmd: MSG_AUTH_STATUS, version: version, body: &AuthenticationStatus{0}}
	client.EnqueueMessage(msg)

	log.Info(client.Connection)

	client.AddClient()

	client.PeerClient.Login()
	//go ChangeUserOnline(uid, 1)
	fmt.Println("tcp login ")
	SetOnline(uid, 1)
	//CountDAU(client.appid, client.uid)
	//atomic.AddInt64(&server_summary.nclients, 1)
}

func (client *Client) AddClient() {
	route := app_route.FindOrAddRoute(client.appid)
	route.AddClient(client)
}

func (client *Client) HandlePing() {
	m := &Message{cmd: MSG_PONG}
	client.EnqueueMessage(m)
	client.timesout=0
	if client.uid == 0 {
		log.Warning("client has't been authenticated")
		return
	}
}

/*func (client *Client) HandleHttp(message *Message) {
	msg := message.body.(*HttpSocket)
	msg.url


	//client.EnqueueMessage(m)

}*/

func (client *Client) HandleACK(ack *MessageACK) {
	log.Info("ack:", ack.seq)
}

//发送等待队列中的消息
func (client *Client) SendMessages(seq int) int {
	var messages *list.List
	client.mutex.Lock()
	if client.messages.Len() == 0 {
		client.mutex.Unlock()
		return seq
	}
	messages = client.messages
	client.messages = list.New()
	client.mutex.Unlock()

	e := messages.Front()
	for e != nil {
		msg := e.Value.(*Message)
		if msg.cmd == MSG_RT || msg.cmd == MSG_IM || msg.cmd == MSG_GROUP_IM {
			atomic.AddInt64(&server_summary.out_message_count, 1)
		}

		if msg.meta != nil {
			seq++
			meta_msg := &Message{cmd: MSG_METADATA, seq: seq, version: client.version, body: msg.meta}
			client.send(meta_msg)
		}
		seq++
		//以当前客户端所用版本号发送消息
		vmsg := &Message{cmd: msg.cmd, seq: seq, version: client.version, flag: msg.flag, body: msg.body}
		client.send(vmsg)

		e = e.Next()
	}
	return seq
}

func (client *Client) Write() {
	seq := 0
	running := true

	//发送在线消息
	for running {
		select {
		case msg := <-client.wt:
			if msg == nil {
				client.close()
				running = false
				log.Infof("client:%d socket closed", client.uid)
				break
			}
			if msg.cmd == MSG_RT || msg.cmd == MSG_IM || msg.cmd == MSG_GROUP_IM {
				//atomic.AddInt64(&server_summary.out_message_count, 1)
			}

			if msg.meta != nil {
				seq++
				meta_msg := &Message{cmd: MSG_METADATA, seq: seq, version: client.version, body: msg.meta}
				client.send(meta_msg)
			}

			seq++
			//以当前客户端所用版本号发送消息
			vmsg := &Message{cmd: msg.cmd, seq: seq, version: client.version, flag: msg.flag, body: msg.body}
			log.Infof("client:%d Write message", client.uid)
			log.Info(vmsg)
			client.send(vmsg)
		case messages := <-client.pwt:
			for _, msg := range messages {
				if msg.cmd == MSG_RT || msg.cmd == MSG_IM || msg.cmd == MSG_GROUP_IM {
					atomic.AddInt64(&server_summary.out_message_count, 1)
				}

				if msg.meta != nil {
					seq++
					meta_msg := &Message{cmd: MSG_METADATA, seq: seq, version: client.version, body: msg.meta}
					client.send(meta_msg)
				}
				seq++
				//以当前客户端所用版本号发送消息
				vmsg := &Message{cmd: msg.cmd, seq: seq, version: client.version, flag: msg.flag, body: msg.body}
				client.send(vmsg)
			}
		case <-client.lwt:
			seq = client.SendMessages(seq)
			log.Info("client.lwt out and it break client id=%d", client.uid)
			break

		}
	}

	//等待200ms,避免发送者阻塞
	t := time.After(200 * time.Millisecond)
	running = true
	for running {
		select {
		case <-t:
			log.Warning("running = false")
			running = false
		case <-client.wt:
			log.Warning("msg is dropped")
		}
	}

	log.Info("write goroutine exit")
}

func (client *Client) Run() {
	go client.CheckAuth()
	go client.Write()
	go client.Read()
	//go client.CheckLive()
}

//检查连接
func (client *Client) CheckAuth() {

	tiker := time.NewTicker(time.Second * 1)
	<-tiker.C
	if !(math.Abs(float64(client.uid)) > 0) {
		log.Info("client close")
		client.close()
	}
}
/*func (client *Client)CheckLive()  {
	tiker := time.NewTicker(time.Second * 2)
	<-tiker.C
	if client.uid==0{
		return
	}
	for{
		tiker := time.NewTicker(time.Second * 15)
		<-tiker.C
		if client.uid==0{
			return
		}
		client.timesout=client.timesout+1
		fmt.Println("检查存活中 client:",client.uid," closed:",client.closed,"timeout:",client.timesout)
		if client.timesout>5{
			client.close()
			//client.HandleClientClosed()
			fmt.Println("检查存活--断开")
			return
		}
	}
}*/

func (client *Client) Handerfunc(pattarn string, handler func(c *Client, r *Request)) {
	if _, exist := client.Handler[pattarn]; exist {
		panic("http: multiple registrations for " + pattarn)
	}
	client.Handler[pattarn] = handler
}

//设置用户在线状态
func SetOnline(uid int64, status int) {

	users.ChangeUserOnline(uid,status)
}
