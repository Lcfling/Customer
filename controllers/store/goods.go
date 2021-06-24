package store

import (
	"encoding/json"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/utils"
)

type AddStore struct {
	controllers.MobileController
}

func (this *AddStore) Get() {
	barcode, _ := this.GetInt64("barcode")
	store_id, _ := this.GetInt64("storeid")
	pro, err := order.GetProductByBigCode(barcode, store_id)
	if err == nil {
		pro.IsBig = 1
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
		this.ServeJSON()
		return
	}
	pro = order.GetProductByStoreCode(barcode, store_id)
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
	this.ServeJSON()
}
func (this *AddStore) Post() {

	//todo 添加入库锁 2秒内 单用户只能入库一次

	id, _ := this.GetInt64("id")
	nums, _ := this.GetInt64("nums")
	if !(nums > 0 && id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "入库数量必须大于0"}
		this.ServeJSON()
		return
	}
	err := order.AddStock(id, nums, this.Uid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "入库错误，请联系管理员"}
		this.ServeJSON()
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
	}
}

//所有分类
type GateTree struct {
	controllers.MobileController
}

func (this *GateTree) Get() {
	tree, err := order.GetTree()
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "empty"}
		this.ServeJSON()
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": tree}
		this.ServeJSON()
	}
}

type SubOrderController struct {
	controllers.UserBaseController
}

func (this *SubOrderController) Post() {
	sublistJson := this.GetString("proList")
	var paytype int64

	paytype = 5 //现金支付
	var sublist []order.SubList

	err := json.Unmarshal([]byte(sublistJson), &sublist)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "json格式错误"}
		this.ServeJSON()
		return
	}

	token := this.GetString("dsn")
	doorinfo, _ := device.GetDiviveByToken(token)
	store_id := doorinfo.StoreId
	if !(store_id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商家信息不存在"}
		this.ServeJSON()
		return
	}
	//生成订单号
	order_id := utils.GetOrderSN()
	enterlog, _ := logs.GetEid(this.Uid)

	err = order.CreatOrder(sublist, store_id, this.Uid, order_id, paytype, enterlog.Id)
	//logs.UpdateOrderid(storeid,this.Uid,order_id)

	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "创建订单错误，请联系客服处理！"}
		this.ServeJSON()
		return
	}
}
