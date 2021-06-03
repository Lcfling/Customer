package users

import (
	"fmt"
	"github.com/Lcfling/Customer/alipay"
	"github.com/Lcfling/Customer/alipay/tools"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/cmd"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/service"
	"github.com/Lcfling/Customer/socket"
	"github.com/astaxie/beego"
	"strconv"
)

type AlipayNotify struct {
	controllers.IndexController
}

func (this *AlipayNotify) Post(){
	data:=this.Input()
	signMap :=make(map[string]string)
	for k,v:= range data {
		signMap[k]=v[0]
	}
	sign:=signMap["sign"]
	sign_type:=signMap["sign_type"]
	delete(signMap, "sign")
	delete(signMap, "sign_type")

	fmt.Println("input data:",signMap,data)


	aliPayPublicKey:=beego.AppConfig.String("ali_publickey")
	publicKey:=tools.FormatPublicKey(aliPayPublicKey)
	signdata:=tools.Signdata(signMap)
	err:=alipay.VerifySign(signdata,sign,sign_type,publicKey)
	if err!=nil{
		this.Ctx.WriteString("faild!")
		fmt.Println("验签失败：",err.Error())
		return
	}


	order_id:=signMap["out_trade_no"]
	out_trade_sn:=signMap["trade_no"]
	moneyf,_:=strconv.ParseFloat(signMap["total_amount"],64)
	money:=int64(moneyf*100)
	orderinfo,err:=order.GetOrderByOrderId(order_id)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": "false", "message": "订单不存在"}
		this.ServeJSON()
		return
	}
	err=order.OrderPaid(order_id,money,out_trade_sn)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": "false", "message": err.Error()}
		this.ServeJSON()
		return
	}
	service_id,_:=service.GetService(orderinfo.StoreId)
	change:=cmd.SendOrderStatus{Cmd:5,Uid:orderinfo.Uid,Storeid:orderinfo.StoreId,Status:1,Ordersn:order_id}
	go socket.SendMessageToPeer(service_id,change)

	this.Ctx.WriteString("success")
	return
}