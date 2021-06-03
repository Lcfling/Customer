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
	errs "errors"
	"time"
)
import log "github.com/golang/glog"

func GetChannel(uid int64) *Channel {
	if uid < 0 {
		uid = -uid
	}
	index := uid % int64(len(route_channels))
	return route_channels[index]
}

func GetRoomChannel(room_id int64) *Channel {
	if room_id < 0 {
		room_id = -room_id
	}
	index := room_id % int64(len(route_channels))
	return route_channels[index]
}

//离线消息推送
func PushMessage(appid int64, uid int64, m *Message) {
	channel := GetChannel(uid)
	channel.Push(appid, []int64{uid}, m)
}

func PublishMessage(appid int64, uid int64, msg *Message) {
	now := time.Now().UnixNano()
	amsg := &AppMessage{appid: appid, receiver: uid, timestamp: now, msg: msg}
	if msg.meta != nil {
		amsg.msgid = msg.meta.sync_key
		amsg.prev_msgid = msg.meta.prev_sync_key
	}
	channel := GetChannel(uid)
	channel.Publish(amsg)
}

func SendAppMessage(appid int64, uid int64, msg *Message) {
	now := time.Now().UnixNano()
	amsg := &AppMessage{appid: appid, receiver: uid, msgid: 0, timestamp: now, msg: msg}
	channel := GetChannel(uid)
	channel.Publish(amsg)
	DispatchMessageToPeer(msg, uid, appid, nil)
}

func DispatchAppMessage(amsg *AppMessage) {
	now := time.Now().UnixNano()
	d := now - amsg.timestamp
	log.Infof("dispatch app message:%s %d %d", Command(amsg.msg.cmd), amsg.msg.flag, d)
	if d > int64(time.Second) {
		log.Warning("dispatch app message slow...")
	}

	if amsg.msgid > 0 {
		if (amsg.msg.flag & MESSAGE_FLAG_PUSH) == 0 {
			log.Fatal("invalid message flag", amsg.msg.flag)
		}
		meta := &Metadata{sync_key: amsg.msgid, prev_sync_key: amsg.prev_msgid}
		amsg.msg.meta = meta
	}
	DispatchMessageToPeer(amsg.msg, amsg.receiver, amsg.appid, nil)
}

func DispatchRoomMessage(amsg *AppMessage) {
	log.Info("dispatch room message", Command(amsg.msg.cmd))
	room_id := amsg.receiver
	DispatchMessageToRoom(amsg.msg, room_id, amsg.appid, nil)
}

func DispatchMessageToPeer(msg *Message, uid int64, appid int64, client *Client) bool {
	route := app_route.FindRoute(appid)
	if route == nil {
		log.Warningf("can't dispatch app message, appid:%d uid:%d cmd:%s", appid, uid, Command(msg.cmd))
		return false
	}
	clients := route.FindClientSet(uid)
	if len(clients) == 0 {
		return false
	}

	for c, _ := range clients {
		if c == client {
			continue
		}
		c.EnqueueNonBlockMessage(msg)
	}
	return true
}

func DispatchMessageToRoom(msg *Message, room_id int64, appid int64, client *Client) bool {
	route := app_route.FindOrAddRoute(appid)
	clients := route.FindRoomClientSet(room_id)

	if len(clients) == 0 {
		return false
	}
	for c, _ := range clients {
		/*if c == client {
			continue
		}*/

		c.EnqueueNonBlockMessage(msg)
	}
	return true
}
func SendMessageToPeer(uid int64, content interface{}) error {
	//
	if !(uid > 0) {
		return errs.New("uid error")
	}
	r, err := json.Marshal(content)

	if err != nil {
		log.Warning("SendMessageToPeer faild:", err)
		return errs.New("json encode error")
	}
	im := &IMMessage{}
	im.sender = 0
	im.receiver = uid
	im.msgid = 0
	im.timestamp = int32(time.Now().Unix())
	im.content = string(r)
	m := &Message{cmd: MSG_IM, version: DEFAULT_VERSION, body: im}
	PublishMessage(1, uid, m)
	return nil
}

//推送给所有用户
func SendMessageToAll(content interface{}) error {
	//

	route := app_route.FindOrAddRoute(1)
	for k, _ := range route.clients {
		uid := k
		if !(uid > 0) {
			return errs.New("uid error")
		}
		r, err := json.Marshal(content)

		if err != nil {
			log.Warning("SendMessageToPeer faild:", err)
			return errs.New("json encode error")
		}
		im := &IMMessage{}
		im.sender = 0
		im.receiver = uid
		im.msgid = 0
		im.timestamp = int32(time.Now().Unix())
		im.content = string(r)

		m := &Message{cmd: MSG_IM, version: DEFAULT_VERSION, body: im}
		PublishMessage(1, uid, m)
	}
	return nil
}
func SendMessageToRoom(room int64, content interface{}) error {

	if !(room > 0) {
		return errs.New("roomid error")
	}
	r, err := json.Marshal(content)

	if err != nil {
		log.Warning("SendMessageToRoom faild:", err)
		return errs.New("json encode error")
	}

	room_im := &RoomMessage{new(RTMessage)}
	room_im.sender = 0
	room_im.receiver = room
	room_im.content = string(r)

	msg := &Message{cmd: MSG_ROOM_IM, body: room_im}

	DispatchMessageToRoom(msg, room, 1, nil)

	amsg := &AppMessage{appid: 1, receiver: room, msg: msg}
	channel := GetRoomChannel(room)
	channel.PublishRoom(amsg)
	return nil
}
