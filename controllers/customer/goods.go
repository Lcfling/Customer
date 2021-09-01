package customer

import (
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/utils"
	"strconv"
)

//需要权限
//二维码查询商品
type GoodsController struct {
	controllers.UserBaseController
}

func (this *GoodsController) Get() {

	barcode, _ := this.GetInt64("barcode")
	token := this.GetString("dsn")

	if !(barcode > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品不存在", "data": ""}
		this.ServeJSON()
		return
	}
	strcode := utils.GetStrPos(strconv.FormatInt(barcode, 10), 0, 7)
	intcode, _ := strconv.ParseInt(strcode, 10, 64)
	strnums := utils.GetStrPos(strconv.FormatInt(barcode, 10), 7, 5)
	intnums, _ := strconv.ParseInt(strnums, 10, 64)

	doorinfo, _ := device.GetDiviveByToken(token)
	store_id := doorinfo.StoreId

	pro, err := order.GetProductByStoreCode(intcode, store_id)
	if err == nil {
		nums := float64(intnums) / float64(pro.Price)
		pro.Nums = nums
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
		this.ServeJSON()
		return
	}

	pro, err = order.GetProductByStoreCode(barcode, store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品不存在", "data": ""}
		this.ServeJSON()
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
		this.ServeJSON()
	}

}
