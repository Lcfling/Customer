package udpsok

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	log "github.com/golang/glog"
	"io"
	"net"
	"sync"
)

type Connection struct {
	conn 	interface{}
	addr 	*net.UDPAddr
	sn 		int64
	lastime int64
	wt 		chan *Message
	rt 		chan []byte

	messages *list.List //待发送的消息队列 FIFO
	mutex    sync.Mutex
}

func (client *Connection) read(data []byte) *Message  {
	//buff := make([]byte, 8)
	flag,cmd,_,sn:=ReadHeader(data)


	message:=new(Message)
	message.cmd=cmd
	message.flag=flag
	message.sn=sn
	if !message.FromData(data){
		log.Warningf("udp message parse error:%d",sn)
	}
	return message
}
func ReadHeader(buff []byte) (int8, int8, int16, int32) {
	var flag int8
	var cmd int8
	var k 	int16
	var sn  int32
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.LittleEndian, &flag)
	binary.Read(buffer, binary.LittleEndian, &cmd)
	binary.Read(buffer, binary.LittleEndian, &k)
	binary.Read(buffer, binary.LittleEndian, &sn)

	return flag,cmd,k,sn
}
func WriteHeader(flag int8, cmd int8, k int16,sn int32 , buffer io.Writer) {
	binary.Write(buffer, binary.LittleEndian, flag)
	binary.Write(buffer, binary.LittleEndian, cmd)
	binary.Write(buffer, binary.LittleEndian, k)
	binary.Write(buffer, binary.LittleEndian, sn)
}
func WriteMessage(w *bytes.Buffer, msg *Message) {
	body := msg.ToData()
	WriteHeader(msg.flag,msg.cmd,msg.k,msg.sn, w)
	w.Write(body)
}
func (client *Connection)SendMessage(msg *Message) error {
	buffer := new(bytes.Buffer)
	WriteMessage(buffer, msg)
	buf := buffer.Bytes()


	_, err:= client.conn.(*net.UDPConn).WriteToUDP(buf, client.addr)
	if err != nil {
		fmt.Println("send message:",buf)
		fmt.Println("发送数据失败!", err)
		return err

	}
	return nil
}
func (client *Connection) close() {
	if conn, ok := client.conn.(*net.UDPConn); ok {
		conn.Close()
	}
}