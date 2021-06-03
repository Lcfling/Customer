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
	"encoding/json"
	log "github.com/golang/glog"
	"strconv"
)
import "sync/atomic"

type RoomClient struct {
	*Connection
	room_id int64
}

func (client *RoomClient) Logout() {
	if client.room_id > 0 {
		channel := GetRoomChannel(client.room_id)
		channel.UnsubscribeRoom(client.appid, client.room_id)
		route := app_route.FindOrAddRoute(client.appid)
		route.RemoveRoomClient(client.room_id, client.Client())
		SendRoomCount(client.room_id, route)
	}
}

func (client *RoomClient) HandleMessage(msg *Message) {
	switch msg.cmd {
	case MSG_ENTER_ROOM:
		client.HandleEnterRoom(msg.body.(*Room))
	case MSG_LEAVE_ROOM:
		client.HandleLeaveRoom(msg.body.(*Room))
	case MSG_ROOM_IM:
		client.HandleRoomIM(msg.body.(*RoomMessage), msg.seq)
	}
}

func (client *RoomClient) HandleEnterRoom(room *Room) {
	if client.uid == 0 {
		log.Warning("client has't been authenticated")
		return
	}

	room_id := room.RoomID()
	log.Info("enter room id:", room_id)
	if room_id == 0 || client.room_id == room_id {
		return
	}
	route := app_route.FindOrAddRoute(client.appid)
	if client.room_id > 0 {
		channel := GetRoomChannel(client.room_id)
		channel.UnsubscribeRoom(client.appid, client.room_id)
		route.RemoveRoomClient(client.room_id, client.Client())
	}

	client.room_id = room_id
	route.AddRoomClient(client.room_id, client.Client())
	SendRoomCount(client.room_id, route)
	channel := GetRoomChannel(client.room_id)
	channel.SubscribeRoom(client.appid, client.room_id)
}

func (client *RoomClient) HandleLeaveRoom(room *Room) {
	if client.uid == 0 {
		log.Warning("client has't been authenticated")
		return
	}

	room_id := room.RoomID()
	log.Info("leave room id:", room_id)
	if room_id == 0 {
		return
	}
	if client.room_id != room_id {
		return
	}

	route := app_route.FindOrAddRoute(client.appid)
	route.RemoveRoomClient(client.room_id, client.Client())
	SendRoomCount(client.room_id, route)
	channel := GetRoomChannel(client.room_id)
	channel.UnsubscribeRoom(client.appid, client.room_id)
	client.room_id = 0
}

func (client *RoomClient) HandleRoomIM(room_im *RoomMessage, seq int) {
	if client.uid == 0 {
		log.Warning("client has't been authenticated")
		return
	}
	room_id := room_im.receiver
	if room_id != client.room_id {
		log.Warningf("room id:%d is't client's room id:%d\n", room_id, client.room_id)
		return
	}

	//判断禁言
	if Isforbiden(room_im.content, client.uid) {
		return
	}

	fb := atomic.LoadInt32(&client.forbidden)
	if fb == 1 {
		log.Infof("room id:%d client:%d, %d is forbidden", room_id, client.appid, client.uid)
		return
	}

	m := &Message{cmd: MSG_ROOM_IM, body: room_im}
	DispatchMessageToRoom(m, room_id, client.appid, client.Client())

	amsg := &AppMessage{appid: client.appid, receiver: room_id, msg: m}
	channel := GetRoomChannel(client.room_id)
	channel.PublishRoom(amsg)

	client.wt <- &Message{cmd: MSG_ACK, body: &MessageACK{seq: int32(seq)}}
}

//发送同步房间人数
func SendRoomCount(roomid int64, route *Route) {
	strings := strconv.FormatInt(roomid, 10)
	maps, _ := GetRedis("roomid_" + strings)

	countmap := make(map[string]int64)
	if maps != "" {
		err := json.Unmarshal([]byte(maps), &countmap)
		if err != nil {
			log.Info("json to map err:", err)
			return
		}
	}

	countmap[macAddr] = int64(len(route.room_clients[roomid]))
	res, err := json.Marshal(countmap)
	if err != nil {
		log.Info("map to json err:", err)
		return
	}
	SetRedis("roomid_"+strings, string(res))
	//err:=json.Marshal(maps)


}
func Isforbiden(content string, uid int64) bool {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(content), &mapResult)
	if err != nil {
		return false
	}
	if _, ok := mapResult["Cmd"]; !ok {
		return false
	}
	if int64(mapResult["Cmd"].(float64)) != 66 {
		return false
	}
	return true
}
