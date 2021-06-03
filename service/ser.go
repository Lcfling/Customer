package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lcfling/Customer/models/cmd"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/socket"
	"github.com/Lcfling/Customer/udpsok"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

//服务代码
var wt chan *logs.EnterLog
var closerStore chan *logs.CloseStore
var doworker chan int64
var store map[int64]int64 //store -> service id
type service map[int64]int64
var services map[int64]service //service id ->

var redis_pool *redis.Pool
func init(){
	wt=make(chan *logs.EnterLog,300)
	closerStore=make(chan *logs.CloseStore,300)
	doworker=make(chan int64,300)
	store=make(map[int64]int64)
	services=make(map[int64]service)
}
func InitData(){
	redis_pool=NewRedisPool(beego.AppConfig.String("redis_host"),beego.AppConfig.String("redis_password"),0)
	//redis_pool=NewRedisPool()
	//utils.DEL(redis_pool,"servicesdata")
	//utils.DEL(redis_pool,"storedata")
	servicesJson,err:=utils.Get(redis_pool,"servicesdata")

	if err==nil {
		json.Unmarshal([]byte(servicesJson.(string)), &services)
	}
	storeJson,err:=utils.Get(redis_pool,"storedata")

	if err==nil {
		json.Unmarshal([]byte(storeJson.(string)), &store)
	}
}
func Run(){
	InitData()
	//redis_pool:=socket.GetRedisHandle()
	for{
		fmt.Println("循环")
		fmt.Println("services:",services)
		fmt.Println("store:",store)
		select {
		case enterlogs:=<-wt:
			_,ok:=store[enterlogs.StoreId]
			if ok {
				//todo 吧用户加入到对应客服的队列下

				// 发送给客服消息 告诉客服某某客户进入

				logs.UpdateServiceId(store[enterlogs.StoreId],enterlogs.Uid)
				fmt.Println("进入房间")
				NotifyService(store[enterlogs.StoreId],enterlogs)
			}else {
				//

				serviceid,err:=GetStartMq()

				fmt.Println("ser:",serviceid," err:",err)
				/*if serviceid==0{
					//todo 当前情况客服为空 请处理
					fmt.Println("todo 当前情况客服为空 请处理")
					continue
				}*/
				if err==nil&&serviceid!=0 {
					//吧客服放入高阶队列
					utils.RPush(redis_pool,"HIGHMQ",serviceid)
					//utils.RPush(redis_pool,"serlist_"+strconv.FormatInt(serviceid.(int64),10),logs.StoreId)
					sid,_:=strconv.ParseInt(serviceid.(string),10,64)
					DataAndServer(sid,enterlogs.StoreId,enterlogs.Uid)
					// 通知客服 店铺进入客户 打开对应店铺的服务窗口
					NotifyService(sid,enterlogs)

				}else{
					// 处理高阶队列情况
					serlist,err:=utils.LRange(redis_pool,"HIGHMQ",0,-1)
					if err!=nil {
						fmt.Println("高阶队列为空")
						continue
					}
					//var index int=0
					var lenhandle int = 1000
					var service_id int64
					service_id=0
					for _,v :=range serlist{
						//客服如果暂时离开状态 则不加入算法中计算

						sid,_:=strconv.ParseInt(string(v.([]uint8)),10,64)

						Info,_:=users.GetUser(sid)
						if Info.Online==0||Info.Isleave==1{
							continue
						}
						if len(services[sid]) <lenhandle{

							//index=k
							service_id=sid
							lenhandle=len(services[sid])
							fmt.Println("len:",lenhandle,"sid",sid)
						}
					}
					//service_id:=serlist[index]
					if service_id==0{
						continue
					}
					fmt.Println("services[service_id][logs.StoreId]",services)

					DataAndServer(service_id,enterlogs.StoreId,enterlogs.Uid)
					// 通知客服 店铺进入客户 打开对应店铺的服务窗口

					NotifyService(service_id,enterlogs)

				}
			}
		case data:=<-closerStore:
			if !(data.StoreId>0) ||!(data.SerivceId>0){
				fmt.Println("监听到错误的关闭店铺参数")
				continue
			}
			delete(services[data.SerivceId],data.StoreId)
			delete(store,data.StoreId)
			servicesJson, _ := json.Marshal(services)
			utils.Set(redis_pool,"servicesdata",string(servicesJson))
			storeJson,_:=json.Marshal(store)
			utils.Set(redis_pool,"storedata",string(storeJson))
			if len(services[data.SerivceId])<1{
				//离开高阶队列
				utils.LRem(redis_pool,"HIGHMQ",0,data.SerivceId)
				//放入低阶队列  判断客服是否离开 是否在线
				Userinfo,_:=users.GetUser(data.SerivceId)
				if Userinfo.Online==1&&Userinfo.Isleave==0{
					utils.RPush(redis_pool,"STARTMQ",data.SerivceId)
				}

			}
			//设置所有当前店铺的用户为离开状态
			logs.CloseStoreById(data.StoreId)
			//发送关灯指令
			CloseTheLight(data.StoreId)

		case Worker:=<-doworker:
			if !(Worker>0){
				fmt.Println("监听到错误的客服id")
				continue
			}
			serlist,err:=utils.LRange(redis_pool,"STARTMQ",0,-1)
			start:=false
			hight:=false
			if err==nil{
				for _,v:=range serlist{
					if string(v.([]uint8))== strconv.FormatInt(Worker,10){
						start=true
					}
				}
			}
			hiserlist,err:=utils.LRange(redis_pool,"HIGHMQ",0,-1)
			if err==nil{
				for _,v:=range hiserlist{
					if string(v.([]uint8))== strconv.FormatInt(Worker,10){
						hight=true
					}
				}
			}

			if !hight && !start {
				fmt.Println("高低队列都没有开始入队")
				Workerstr:=strconv.FormatInt(Worker,10)
				utils.RPush(redis_pool,"STARTMQ",Workerstr)
				k,_:=utils.LRange(redis_pool,"STARTMQ",0,-1)
				fmt.Println("------k:",k)


			}else {
				fmt.Println("已经在队列中",hight,start)
			}
		}
	}
}

func GetChannel() chan *logs.EnterLog {
	return wt
}
func GetCloseChannel() chan *logs.CloseStore {
	return closerStore
}
func GetWorking() chan int64{
	return doworker
}

func GetStartMq() (interface{},error) {
	serviceid,err:=utils.LPop(redis_pool,"STARTMQ")
	fmt.Println("serivceid ",serviceid,"err:",err)
	if err!=nil{
		return 0,err
	}
	//判断客服是否处于工作状态中
	uid,_:=strconv.ParseInt(serviceid.(string),10,64)
	Info,_:=users.GetUser(uid)

	if Info.Online==1&&Info.Isleave==0{
		return serviceid,nil
	}else{
		fmt.Println("二次循环了")
		return GetStartMq()
	}

}

func DataAndServer(serviceid int64,store_id int64,uid int64){
	logs.UpdateServiceId(serviceid,uid)
	serhandle,ok:=services[serviceid]
	if !ok {
		serhandle=make(service)
	}
	serhandle[store_id]=store_id
	services[serviceid]=serhandle
	store[store_id]=serviceid
	if _,ok:=services[0];ok{
		delete(services,0)
		delete(store,0)
		delete(services,1)
		delete(store,1)

	}

	//var servicesData map[string]interface{}
	servicesJson, _ := json.Marshal(services)
	utils.Set(redis_pool,"servicesdata",string(servicesJson))
	storeJson,_:=json.Marshal(store)
	utils.Set(redis_pool,"storedata",string(storeJson))
}
func NotifyService(serviceid int64,data *logs.EnterLog){


	isopendoor:=data.LeaveTime
	data.LeaveTime=0
	data.ServiceId=serviceid
	logs.AddEnterLog(*data)
	//通知代码
	// 开门
	doorinfo,err:=device.GetDiviveById(data.DoorId)

	if err!=nil{
		fmt.Println("device.GetDiviveById",err.Error())
	}
	if isopendoor==1{
		err=udpsok.HandleOpenDoor(doorinfo.Devicesn,doorinfo.Nums)
		if err!=nil{
			fmt.Println("udpsok.HandleOpenDoor",err.Error())
		}
	}
	go func() {
		time.Sleep(time.Duration(1)*time.Second)
		OpenTheLight(data.StoreId)
	}()
	Enter:=cmd.UserEnter{Cmd:1,Uid:data.Uid,Storeid:data.StoreId,Doorid:data.DoorId}
	fmt.Println("发送消息",Enter)
	go socket.SendMessageToPeer(serviceid,Enter)
}
//判断客服是否服务于店铺
func IsServiceStore(service_id int64,store_id int64) bool{
	_,ok:=store[store_id]
	if ok&&store[store_id]==service_id{
		return true
	}
	return false
}
func GetService(store_id int64) (int64,error ){
	_,ok:=store[store_id]
	if ok{
		return store[store_id],nil
	}else {
		return 0,errors.New("未找到店铺对应的客服")
	}
}
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

func GetStoreList(service_id int64) service{
	list,ok:=services[service_id]
	if !ok{
		return nil
	}else{
		return list
	}
}
func CloseTheLight(store_id int64){
	door,_:=device.GetLampByStore(store_id,3)
	udpsok.OpenDoor(door.Devicesn,door.Nums, func(e bool) {
		fmt.Println("关灯指令状态：",e)
	})
}
func OpenTheLight(store_id int64){
	door,_:=device.GetLampByStore(store_id,2)
	udpsok.OpenDoor(door.Devicesn,door.Nums, func(e bool) {
		fmt.Println("开灯指令状态：",e)
	})
}