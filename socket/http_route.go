package socket

import (
	log "github.com/golang/glog"
	"net/http"
)

func SetHttpRoute() {
	http.HandleFunc("/summary", Summary)
	http.HandleFunc("/serverping", ping)

	//post 推送消息
	http.HandleFunc("/postpeermessage", PostPeerMessage)
	http.HandleFunc("/postmessagetoall", PostMessageToAll)

}

//web服务
func StartHttpServer(addr string) {
	//rpc function
	//http.HandleFunc("/post_group_notification", PostGroupNotification)

	//请在此处添加http路由规则

	handler := loggingHandler{http.DefaultServeMux}
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal("http server err:", err)
	}
}

//web服务 https
func StartHttpServerSSL(addr string, crt string, key string) {
	//rpc function
	//http.HandleFunc("/post_group_notification", PostGroupNotification)

	//请在此处添加http路由规则

	handler := loggingHandler{http.DefaultServeMux}
	err := http.ListenAndServeTLS(addr, crt, key, handler)
	if err != nil {
		log.Fatal("http server err:", err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	//WriteHttpObj

	b := "success"
	w.Write([]byte(b))
}
