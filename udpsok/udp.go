package udpsok

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)
var mutex    sync.Mutex
var Clients map[int64]*Client

var Tests map[int64]string
func Run()  {
	// 创建监听
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 61005,
	})
	if err != nil {
		fmt.Println("监听失败!", err)
		return
	}
	defer socket.Close()
	Clients=make(map[int64]*Client)
	//Tests[32]="nihao"
	for {
		// 读取数据
		data := make([]byte, 64)
		_, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("读取数据失败!", err)
			continue
		}

		go HandelConn(socket,data,remoteAddr)


		// 发送数据

		/*for{
			senddata := []byte("hello client!")
			_, err = socket.WriteToUDP(senddata, remoteAddr)
			if err != nil {
				return
				fmt.Println("发送数据失败!", err)
			}
		}*/

	}
}
func HandelConn(conn *net.UDPConn,data []byte,addr *net.UDPAddr)  {
	var flag int8
	var cmd int8
	var k int16
	var sn32 int32
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.LittleEndian, &flag)
	binary.Read(buffer, binary.LittleEndian, &cmd)
	binary.Read(buffer, binary.LittleEndian, &k)
	binary.Read(buffer, binary.LittleEndian, &sn32)
	sn:=int64(sn32)

	c,ok:=Clients[sn]
	if !ok{
		fmt.Println("新建连接!sn",sn,ok)
		mutex.Lock()
		Clients[sn]=NewClient(conn,sn,addr)
		mutex.Unlock()
	}else{
		Clients[sn].UpdateAddr(addr)
		c.rt <- data

	}

}

func OpenDoor(sn int64,num int8 ,f func(e bool)) {


	msg:=&Message{flag:23,cmd:MSG_OPEN,k:0,sn:int32(sn),body:&OpenMessage{index:num}}
	c,ok:=Clients[sn]
	if !ok{
		f(false)
		return
	}
	c.wt<-msg
	c.f=f
}
func HandleOpenDoor(sn int64,num int8) error {
	var channel chan bool
	channel= make(chan bool,20)
	OpenDoor(sn,num, func(e bool) {
		channel <- e
	})
	var wait  chan bool
	wait= make(chan bool,20)
	go func() {
		time.Sleep(time.Duration(10)*time.Second)
		wait<-false
	}()
	running:=true
	for  running{
		select {
		case status:= <- channel:
			running=false
			if status{
				return nil
			}else{
				return errors.New("请联系管理员")
			}

		case <-wait:
			running=false
			return errors.New("响应超时")
		}
	}
	return nil
}