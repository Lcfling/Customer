package socket

import (
	errs "errors"
	log "github.com/golang/glog"
	"net/http"
	"strconv"
	"time"
)
import "encoding/json"
import "os"
import "runtime"
import "runtime/pprof"

type ServerSummary struct {
	nconnections      int64
	nclients          int64
	in_message_count  int64
	out_message_count int64
}

func NewServerSummary() *ServerSummary {
	s := new(ServerSummary)
	return s
}

func Summary(rw http.ResponseWriter, req *http.Request) {
	obj := make(map[string]interface{})
	obj["goroutine_count"] = runtime.NumGoroutine()
	obj["connection_count"] = server_summary.nconnections
	obj["client_count"] = server_summary.nclients
	obj["in_message_count"] = server_summary.in_message_count
	obj["out_message_count"] = server_summary.out_message_count

	res, err := json.Marshal(obj)
	if err != nil {
		log.Info("json marshal:", err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	_, err = rw.Write(res)
	if err != nil {
		log.Info("write err:", err)
	}
	return
}

func Stack(rw http.ResponseWriter, req *http.Request) {
	pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	rw.WriteHeader(200)
}

func WriteHttpError(status int, err string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	meta := make(map[string]interface{})
	meta["code"] = status
	meta["message"] = err
	obj["meta"] = meta
	b, _ := json.Marshal(obj)
	w.WriteHeader(status)
	w.Write(b)
}

func WriteHttpFailed(info string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	meta := make(map[string]interface{})
	meta["status"] = 0
	meta["info"] = info
	meta["data"] = nil
	b, _ := json.Marshal(meta)
	w.WriteHeader(200)
	w.Write(b)
}

func WriteHttpObj(data map[string]interface{}, info string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	obj["data"] = data
	obj["info"] = info
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func WriteHttpFile(data map[string]interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(data)
	w.Write(b)
}
func WriteHttpOutLogin(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	obj["data"] = nil
	obj["info"] = "登录状态验证失败"
	obj["status"] = 2
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func WriteHttpEmpty(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	obj["data"] = nil
	obj["info"] = "空列表数据"
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func WriteHttpArray(data []map[string]interface{}, info string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	obj["data"] = data
	obj["info"] = info
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func WriteHttpArrays(data []map[string]map[string]interface{}, info string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	obj["data"] = data
	obj["info"] = info
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func PostPeerMessage(w http.ResponseWriter, r *http.Request) {
	uid, _ := strconv.Atoi(r.PostFormValue("uid"))
	appid, _ := strconv.Atoi(r.PostFormValue("appid"))
	if !(uid > 0) {
		log.Info("uid error")
		WriteHttpFailed("uid error", w)
	}
	if !(appid > 0) {
		log.Info("appid error")
		WriteHttpFailed("appid error", w)
	}
	content := r.PostFormValue("content")

	im := &IMMessage{}
	im.sender = 0
	im.receiver = int64(uid)
	im.msgid = 0
	im.timestamp = int32(time.Now().Unix())
	im.content = content

	m := &Message{cmd: MSG_IM, version: DEFAULT_VERSION, body: im}
	PublishMessage(int64(appid), int64(uid), m)
	WriteHttpObj(nil, "success", w)
	//return nil
}
func PostMessageToAll(w http.ResponseWriter, r *http.Request) {

	appid, _ := strconv.Atoi(r.PostFormValue("appid"))
	if !(appid > 0) {
		log.Info("appid error")
	}
	content := r.PostFormValue("content")

	route := app_route.FindOrAddRoute(int64(appid))
	for k, _ := range route.clients {
		uid := k
		if !(uid > 0) {
			log.Info(errs.New("uid error"))
		}
		im := &IMMessage{}
		im.sender = 0
		im.receiver = uid
		im.msgid = 0
		im.timestamp = int32(time.Now().Unix())
		im.content = content

		m := &Message{cmd: MSG_IM, version: DEFAULT_VERSION, body: im}
		PublishMessage(int64(appid), uid, m)
	}
}
