package customer

import (
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/cmd"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/service"
	"github.com/Lcfling/Customer/socket"
	"github.com/Lcfling/Customer/udpsok"
	"github.com/astaxie/beego"
	"time"
)
//需要权限
//客户进入购物场景
type EnterController struct {
	controllers.UserBaseController
}

func (this *EnterController) Post() {

	token:=this.GetString("dsn")
	doorinfo,_:=device.GetDiviveByToken(token)
	store_id:=doorinfo.StoreId
	door_id:=doorinfo.Id
	if !(store_id>0){
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误的门店，请重新扫码开门"}
		this.ServeJSON()
		return
	}
	_,err:=device.GetDiviveById(door_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "未查询到相关门禁"}
		this.ServeJSON()
		return
	}
	storeinfo,err:=order.GetStoreById(store_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "此店铺歇业中"}
		this.ServeJSON()
		return
	}
	if storeinfo.Closed==1{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "此店铺歇业中"}
		this.ServeJSON()
		return
	}
	// 入队客服

	o,_:=this.GetInt64("o")



	//var In controllers.UserInfoIn
	In:=new(logs.EnterLog)
	In.Uid=this.Uid
	In.StoreId=store_id
	In.DoorId=door_id
	In.EnterTime=time.Now().Unix()
	In.Status=0
	// 入队
	if o==1{
		In.LeaveTime=1
	}



	service.GetChannel()<-In


	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
	}


}


//客户离开场景
type LeaveController struct {
	controllers.UserBaseController
}
func (this *LeaveController) Post() {
	token:=this.GetString("dsn")
	doorinfo,_:=device.GetDiviveByToken(token)
	store_id:=doorinfo.StoreId
	door_id:=doorinfo.Id
	if !(store_id>0){
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误的门店，请重新扫码开门"}
		this.ServeJSON()
		return
	}
	_,err:=device.GetDiviveById(door_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "未查询到相关门禁"}
		this.ServeJSON()
		return
	}
	storeinfo,err:=order.GetStoreById(store_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "此店铺歇业中"}
		this.ServeJSON()
		return
	}
	if storeinfo.Closed==1{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "此店铺歇业中"}
		this.ServeJSON()
		return
	}


	//更新离开信息

	logs.CustomerLeave(this.Uid)
	//发送离开信息给客服
	serviceid,err:=service.GetService(store_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
		return
	}
	leave:=cmd.OpenDoorInfo{Cmd:2,Uid:this.Uid,Storeid:store_id,Doorid:door_id}
	go socket.SendMessageToPeer(serviceid,leave)



	err=udpsok.HandleOpenDoor(doorinfo.Devicesn,doorinfo.Nums)
	if err==nil{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "祝您生活愉快"}
		this.ServeJSON()
		return
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "请重新尝试，或者联系客服！"}
		this.ServeJSON()
		return
	}



}

type StoreInfo struct {
	controllers.UserBaseController
}

func (this *StoreInfo) Get() {
	token:=this.GetString("dsn")
	doorinfo,_:=device.GetDiviveByToken(token)
	store_id:=doorinfo.StoreId
	storeInfo,err:=order.GetStoreById(store_id)
	if err!=nil {
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "店铺信息","data":storeInfo}
		this.ServeJSON()
	}
}


/*type Webview struct {
	controllers.IndexController
}

func (this *Webview) Get(){
	token:=this.GetString("dsn")
	doorinfo,_:=device.GetDiviveByToken(token)
	store_id:=doorinfo.StoreId
	door_id:=doorinfo.Id


}*/

type GetPhone struct {
	controllers.IndexController
}

func (this *GetPhone) Get() {
	/*ciphertext:=this.GetString("ciphertext")
	if this.UserAgent=="weixin"{

	}else{

	}*/

	phone:=beego.AppConfig.String("service_phone")
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "店铺不存在","data":phone}
	this.ServeJSON()
	return
}