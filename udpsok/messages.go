package udpsok

import (
	"bytes"
	"encoding/binary"
)

type MessageCreator func() IMessage
var message_creators map[int]MessageCreator = make(map[int]MessageCreator)

const MSG_PING  = 32
const MSG_OPEN  = 64

func init()  {
	message_creators[MSG_PING] = func() IMessage { return new(PingMessage) }
	message_creators[MSG_OPEN] = func() IMessage { return new(OpenMessage) }
}

type IMessage interface {
	ToData() []byte
	FromData(buff []byte) bool
}


type Message struct {
	flag int8
	cmd  int8
	k 	int16
	sn  int32
	body interface{}
}
func (message *Message) ToData() []byte {
	if message.body != nil {
		if m, ok := message.body.(IMessage); ok {
			return m.ToData()
		}
		return nil
	} else {
		return nil
	}
}
func (message *Message) FromData(buff []byte) bool {
	cmd := message.cmd
	if creator, ok := message_creators[int(cmd)]; ok {

		c := creator()
		r := c.FromData(buff)

		message.body = c
		return r
	}
	return len(buff) == 0
}
type PingMessage struct {
	index    int32
	remark  int8
	active  int8
	doornumber     int8
	inout int8//1进门 2 出门
	cardnums int32
	lasttime string //7为 len=7
	result int8
	doorstate1 int8
	doorstate2 int8
	doorstate3 int8
	doorstate4 int8
	doorbut1 int8
	doorbut2 int8
	doorbut3 int8
	doorbut4 int8
	fault int8 //故障  0 无故障 1 有故障
	nowtime string//3位
	snumber int32//流水 4位
}
func (ping *PingMessage) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, ping.index)

	buf := buffer.Bytes()
	return buf
}

func (ping *PingMessage) FromData(buff []byte) bool {
	buffer := bytes.NewBuffer(buff)
	wast:=make([]byte,8)
	binary.Read(buffer, binary.LittleEndian, &wast)
	binary.Read(buffer, binary.LittleEndian, &ping.index)
	binary.Read(buffer, binary.LittleEndian, &ping.remark)
	binary.Read(buffer, binary.LittleEndian, &ping.active)
	binary.Read(buffer, binary.LittleEndian, &ping.doornumber)
	binary.Read(buffer, binary.LittleEndian, &ping.inout)
	binary.Read(buffer, binary.LittleEndian, &ping.cardnums)
	lasttime := make([]byte, 7)
	binary.Read(buffer, binary.LittleEndian, &lasttime)
	binary.Read(buffer, binary.LittleEndian, &ping.result)
	binary.Read(buffer, binary.LittleEndian, &ping.doorstate1)
	binary.Read(buffer, binary.LittleEndian, &ping.doorstate2)
	binary.Read(buffer, binary.LittleEndian, &ping.doorstate3)
	binary.Read(buffer, binary.LittleEndian, &ping.doorstate4)
	binary.Read(buffer, binary.LittleEndian, &ping.doorbut1)
	binary.Read(buffer, binary.LittleEndian, &ping.doorbut2)
	binary.Read(buffer, binary.LittleEndian, &ping.doorbut3)
	binary.Read(buffer, binary.LittleEndian, &ping.doorbut4)
	binary.Read(buffer, binary.LittleEndian, &ping.fault)
	nowtime:=make([]byte,3)
	binary.Read(buffer, binary.LittleEndian, &nowtime)
	binary.Read(buffer, binary.LittleEndian, &ping.snumber)
	return true
}

type OpenMessage struct {
	index int8 //门序号  或者操作状态
}
func (open *OpenMessage) ToData() []byte {
	data:=make([]byte,56)
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, open.index)
	buf := buffer.Bytes()
	for v,k:=range buf{
		data[v]=k
	}
	return data
}

func (open *OpenMessage) FromData(buff []byte) bool {
	buffer := bytes.NewBuffer(buff)
	wast:=make([]byte,8)
	binary.Read(buffer, binary.LittleEndian, &wast)
	binary.Read(buffer, binary.LittleEndian, &open.index)
	return true
}

