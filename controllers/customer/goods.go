package customer

import (
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/order"
)

//需要权限
//二维码查询商品
type GoodsController struct {
	controllers.UserBaseController
}

func (this *GoodsController) Get() {

	barcode,_:=this.GetInt64("barcode")
	token:=this.GetString("dsn")
	doorinfo,_:=device.GetDiviveByToken(token)
	store_id:=doorinfo.StoreId
	pro:=order.GetProductByStoreCode(barcode,store_id)

	this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":pro}
	this.ServeJSON()
}
