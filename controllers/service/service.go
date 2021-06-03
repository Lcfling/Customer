package service

import (
	"fmt"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/service"
	"github.com/Lcfling/Customer/udpsok"
	"github.com/Lcfling/Customer/utils"
	"strconv"
)

//关闭当前门店服务
type CloserController struct {
	controllers.BaseController
}
func (this *CloserController)Post(){
	store_id,_:=this.GetInt64("storeid")

	if !(store_id>0){
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误的店铺id"}
		this.ServeJSON()
		return
	}

	if !service.IsServiceStore(this.Uid,store_id){
		fmt.Println("stid",store_id)
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "该店铺不在您的队列中，无权限"}
		this.ServeJSON()
		return
	}

	closer:=new(logs.CloseStore)
	closer.StoreId=store_id
	closer.SerivceId=this.Uid
	service.GetCloseChannel()<-closer
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "服务关闭成功"}
	this.ServeJSON()
}

type OpenDoor struct {
	controllers.BaseController
}

func (this *OpenDoor) Post() {
	door_id,_:=this.GetInt64("doorid")
	doorinfo,_:=device.GetDiviveById(door_id)

	err:=udpsok.HandleOpenDoor(doorinfo.Devicesn,doorinfo.Nums)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
	}
}

//获取设备
type Device struct {
	controllers.BaseController
}

func (this *Device)Get(){
	door_id,_:=this.GetInt64("doorid")
	doorInfo,err:=device.GetDiviveById(door_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":doorInfo}
		this.ServeJSON()
	}
}

type GetVideo struct {
	controllers.BaseController
}

func (this *GetVideo) Get(){
	store_id,_:=this.GetInt64("storeid")
	Videos,err:=device.GetVideo(store_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":Videos}
		this.ServeJSON()
	}
}

type LeaveMoment struct {
	controllers.BaseController
}

func (this *LeaveMoment) Post() {
	info,err:=users.GetUser(this.Uid)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
	}
	if info.Isleave==0{
		users.UpdateLeave(this.Uid,1)
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "状态更改成功,当前状态为离开状态","data":1}
		this.ServeJSON()
	}else{
		users.UpdateLeave(this.Uid,0)
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "状态更改成功,当前状态为工作状态","data":0}
		this.ServeJSON()
	}
}


//退出登录
type LoginOut struct {
	controllers.BaseController
}
func (this *LoginOut) Post() {
	serlist:=service.GetStoreList(this.Uid)
	if !(serlist==nil || len(serlist)==0){
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "有未结束服务的店铺，无法退出"}
		this.ServeJSON()
		return
	}
	users.UpdateLeave(this.Uid,1)
	users.UpdateOnline(this.Uid,0)
	mark:="客服登出"
	go logs.AddUserlog(this.Uid,mark,2,this.UserType)
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "success"}
	this.ServeJSON()
	return
}
type UserInfo struct {
	controllers.BaseController
}

func (this *UserInfo) Get() {
	Uid,_:=this.GetInt64("uid")
	user,err:=users.GetUser(Uid)

	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "用户不存在"}
		this.ServeJSON()
	}else{
		user.Token=""
		user.Pwd=""
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":user}
		this.ServeJSON()
	}
}

type GetOrderDetail struct {
	controllers.BaseController
}
func (this *GetOrderDetail)Get(){
	ordersn:=this.GetString("ordersn")

	_,list,err:=order.ListByOrder(ordersn)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "未找到相应订单"}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":list}
		this.ServeJSON()
	}
}
type WorkingStart struct {
	controllers.BaseController
}

func (this *WorkingStart) Post() {
	uid:=this.Uid
	users.UpdateLeave(uid,0)
	service.GetWorking()<-uid
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "success"}
	this.ServeJSON()
}
type IntData struct {
	controllers.BaseController
}

func (this *IntData) Get() {
	uid:=this.Uid
	data:=service.GetStoreList(uid)

	serlist,err:=utils.LRange(models.GetRedis(),"STARTMQ",0,-1)
	start:=false
	hight:=false
	if err==nil{
		for _,v:=range serlist{
			if string(v.([]uint8))== strconv.FormatInt(this.Uid,10){
				start=true
			}
		}
	}
	hiserlist,err:=utils.LRange(models.GetRedis(),"HIGHMQ",0,-1)
	if err==nil{
		for _,v:=range hiserlist{
			if string(v.([]uint8))== strconv.FormatInt(this.Uid,10){
				hight=true
			}
		}
	}

	var status int
	if !hight && !start {
		status=0;
	}else {
		status=1;
	}
	if len(data)>0{
		for k,_:=range data  {
			storeInfo,_:=order.GetStoreById(data[k])
			data[k]=storeInfo.Closed
		}
	}
	res:=map[string]interface{}{"list":data,"status":status}
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":res}
	this.ServeJSON()
}

type GetAccessToken struct {
	controllers.BaseController
}

func (this *GetAccessToken)Get(){
	id,_:=this.GetInt64("id")

	VideoInfo,err:=device.GetVideoById(id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "获取监控失败"}
		this.ServeJSON()
		return
	}
	if VideoInfo.Appid==""||VideoInfo.Secret==""{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "请联系管理员完善监控信息"}
		this.ServeJSON()
		return
	}


	accesstoken,exp,err:=utils.YSGetAccesstoken(VideoInfo.Appid,VideoInfo.Secret)
	exps:=strconv.FormatInt(exp,10)
	device.UpdateVideoToken(accesstoken,exps,VideoInfo.Appid)

	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":map[string]interface{}{"accesstoken":accesstoken,"exp":exp}}
		this.ServeJSON()
	}
}

type StoreInfo struct {
	controllers.BaseController
}
func (this *StoreInfo)Get(){
	store_id,_:=this.GetInt64("storeid")
	storeInfo,err:=order.GetStoreById(store_id)
	if err!=nil {
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":storeInfo}
		this.ServeJSON()
	}
}


//上报信息
type ReportInfo struct {
	controllers.BaseController
}

func (this *ReportInfo)Post(){
	mark:=this.GetString("mark")
	cid,_:=this.GetInt64("uid")
	order_id:=this.GetString("orderid")
	store_id,_:=this.GetInt64("storeid")

	//store_id

	r:=new(order.Report)

	r.Mark=mark
	r.Cid=cid
	r.StoreId=store_id
	r.OrderId=order_id
	r.Uid=this.Uid

	id,err:=order.SaveReport(r)
	if err!=nil {
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "上报失败，请联系管理员"}
		this.ServeJSON()
		return
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":id}
		this.ServeJSON()
	}
}

type StoreStatus struct {
	controllers.BaseController
}
func (this *StoreStatus)Post(){
	store_id,_:=this.GetInt64("storeid")
	if store_id==0{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	}
	status,err:=order.UpdateClosed(store_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "店铺不存在,404"}
		this.ServeJSON()
		return
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":status}
		this.ServeJSON()
		return
	}
}