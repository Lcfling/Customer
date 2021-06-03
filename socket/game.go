package socket

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
	"net/http"
	"time"
)
import log "github.com/golang/glog"

var (
	VERSION       string
	BUILD_TIME    string
	GO_VERSION    string
	GIT_COMMIT_ID string
	GIT_BRANCH    string
)

var config *Config
var redis_pool *redis.Pool //redis全局句柄
//route server
var route_channels []*Channel
var app_route *AppRoute
var macAddr string //服务器mac地址
var server_summary *ServerSummary
var ChangeUserOnline func(int64,int)
func init() {
	app_route = NewAppRoute()
	server_summary = NewServerSummary()
}

type loggingHandler struct {
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*") //header 变量
	if r.Method != "OPTIONS" {
		if r.Method != "POST" {
			WriteHttpError(400, "请求错误", w)
			return
		}
	}
	log.Infof("http request:%s %s %s header:%s", r.RemoteAddr, r.Method, r.URL, r.Header)
	h.handler.ServeHTTP(w, r)
}

//redis句柄函数
func NewRedisPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   500,
		IdleTimeout: 480 * time.Second,
		Dial: func() (redis.Conn, error) {
			timeout := time.Duration(2) * time.Second
			c, err := redis.DialTimeout("tcp", server, timeout, 0, 0)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if db > 0 && db < 16 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}

func Run(f func(int64,int)) {

	fmt.Printf("Version:     %s\nBuilt:       %s\nGo version:  %s\nGit branch:  %s\nGit commit:  %s\n", VERSION, BUILD_TIME, GO_VERSION, GIT_BRANCH, GIT_COMMIT_ID)


	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("usage: im config")
		return
	}
	macAddr = getMacAddrs()
	config = read_cfg(flag.Args()[0])
	log.Infof("net macAddr:%d\n", macAddr)
	log.Infof("port:%d\n", config.port)

	log.Infof("redis address:%s password:%s db:%d\n",
		config.redis_address, config.redis_password, config.redis_db)

	log.Info("storage addresses:", config.storage_rpc_addrs)
	log.Info("route addressed:", config.route_addrs)
	log.Info("group route addressed:", config.group_route_addrs)
	log.Info("kefu appid:", config.kefu_appid)
	log.Info("pending root:", config.pending_root)

	log.Infof("ws address:%s wss address:%s", config.ws_address, config.wss_address)
	log.Infof("cert file:%s key file:%s", config.cert_file, config.key_file)

	redis_pool = NewRedisPool(config.redis_address, config.redis_password,
		config.redis_db)

	fmt.Println("config:",config)
	ChangeUserOnline=f//更改用户在线状态

	//路由频道 配置集群使用
	route_channels = make([]*Channel, 0)
	for _, addr := range config.route_addrs {
		channel := NewChannel(addr, DispatchAppMessage, DispatchGroupMessage, DispatchRoomMessage)
		channel.Start()
		route_channels = append(route_channels, channel)
	}

	//redis订阅 禁烟消息
	go ListenRedis()
	//go SyncKeyService()

	/*SetHttpRoute() //http 路由
	go StartHttpServer(config.http_listen_address)*/
	if len(config.https_listen_address) > 0 && len(config.cert_file) > 0 && len(config.key_file) > 0 {
		go StartHttpServerSSL(config.https_listen_address, config.cert_file, config.key_file)
	}

	//StartRPCServer(config.rpc_listen_address)

	if len(config.ws_address) > 0 {
		log.Info("run wsserver self:")
		go StartWSServer(config.ws_address)
	}
	if len(config.wss_address) > 0 && len(config.cert_file) > 0 && len(config.key_file) > 0 {
		go StartWSSServer(config.wss_address, config.cert_file, config.key_file)
	}

	if config.ssl_port > 0 && len(config.cert_file) > 0 && len(config.key_file) > 0 {
		go ListenSSL(config.ssl_port, config.cert_file, config.key_file)
	}

	log.Info("all running")
	ListenClient()
	log.Infof("exit")

}


func DispatchGroupMessage(amsg *AppMessage) {

}

func getMacAddrs() (macAddr string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return macAddr
	}

	for _, netInterface := range netInterfaces {
		macAddr = netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		} else {
			break
		}

	}
	return macAddr
}

func GetRedisHandle() *redis.Pool {
	return redis_pool
}
//定时创建表任务  八小时循环
/*func taskCreatetable() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("taskCreatetable产生的异常：", err) // 这里的err其实就是panic传入的内容
		}
	}()
	tiker := time.NewTicker(time.Second * 8 * 60 * 60)
	//<-tiker.C
	for i := 0; i > -1; i++ {
		db := Mysql()
		model.CreateTableOrder(db)
		model.CreateTableGameRecord(db)
		model.CreateTableUserBillflow(db)
		model.CreateTableUserFee(db)
		model.DeleteUserLoginLog7(db)
		db.Close()
		<-tiker.C

	}
}*/
