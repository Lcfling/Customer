package store

import (
	"encoding/json"
	"fmt"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/utils"
	"html"
	"strconv"
	"strings"
)

type AddStock struct {
	controllers.MobileController
}

func (this *AddStock) Get() {
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
func (this *AddStock) Post() {

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
	controllers.StaffController
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
	controllers.StaffController
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

	store_id, _ := this.GetInt64("storeid")
	storeInfo, err := order.GetStoreById(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择正确的店铺"}
		this.ServeJSON()
		return
	}
	if storeInfo.Uid != this.Uid {
		if !this.IsStorePower(store_id) {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
			this.ServeJSON()
			return
		}
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
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": order_id}
		this.ServeJSON()
		return
	}
}

type OrderPaid struct {
	controllers.StaffController
}

func (this *OrderPaid) Get() {
	order_id := this.GetString("orderid")
	m, _ := this.GetInt64("money")
	err := order.OrderPaidUnderLine(order_id, m, "0")
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
	}
}

//添加商品
type ProductAddController struct {
	controllers.StaffController
}

func (this *ProductAddController) Get() {
	barcode, _ := this.GetInt64("barcode")
	pro := order.GetProductByCode(barcode)
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
	this.ServeJSON()
}
func (this *ProductAddController) Post() {
	store_id, _ := this.GetInt64("storeid")
	images := html.EscapeString(this.GetString("images"))
	proname := html.EscapeString(this.GetString("proname"))
	proinfo := html.EscapeString(this.GetString("proinfo"))
	keyword := html.EscapeString(this.GetString("keyword"))
	big_code, _ := this.GetInt64("bigcode")
	exchange, _ := this.GetInt64("exchange")
	bar_code, _ := this.GetInt64("barcode")
	cate_id, _ := this.GetInt64("cate_id")
	cost, _ := this.GetInt64("cost")
	price, _ := this.GetInt64("price")
	unit_name := html.EscapeString(this.GetString("unit_name"))
	sort, _ := this.GetInt64("sort")

	if bar_code == 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "条形码不能空"}
		this.ServeJSON()
		return
	}
	if price == 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "价格不能空"}
		this.ServeJSON()
		return
	}
	if !(store_id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺id错误"}
		this.ServeJSON()
		return
	}
	storeInfo, err := order.GetStoreById(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择正确的店铺"}
		this.ServeJSON()
		return
	}
	if storeInfo.Uid != this.Uid {
		if !this.IsStorePower(store_id) {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
			this.ServeJSON()
			return
		}
	}

	var pro order.Product
	pro.StoreId = store_id
	pro.Image = images
	pro.ProName = proname
	pro.ProInfo = proinfo
	pro.Keyword = keyword
	pro.Bigcode = strconv.FormatInt(big_code, 10)
	pro.Exchange = exchange
	pro.BarCode = strconv.FormatInt(bar_code, 10)
	pro.CateId = cate_id
	pro.Cost = cost
	pro.Price = price
	pro.UnitName = unit_name
	pro.Sort = sort
	id, err := order.ProductAdd(pro)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "入库信息成功", "data": fmt.Sprintf("%d", id)}
	} else {
		if strings.Contains(err.Error(), "Error 1062") {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品入库重复"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "入库信息失败：" + err.Error()}
		}

	}
	this.ServeJSON()
}

type ProductEdit struct {
	controllers.StaffController
}

func (this *ProductEdit) Get() {
	id, _ := this.GetInt64("id")
	//store_id,_:=this.GetInt64("storeid")
	if !(id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品id错误"}
		this.ServeJSON()
		return
	}

	pro, err := order.GetProductById(id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品不存在"}
		this.ServeJSON()
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
		this.ServeJSON()
	}

}
func (this *ProductEdit) Post() {
	store_id, _ := this.GetInt64("storeid")
	id, _ := this.GetInt64("id")
	if !(id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品id错误"}
		this.ServeJSON()
		return
	}
	if !(store_id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺id错误"}
		this.ServeJSON()
		return
	}
	storeInfo, err := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		if !this.IsStorePower(store_id) {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
			this.ServeJSON()
			return
		}
	}

	var p order.Product

	proname := this.GetString("proname")
	unitname := this.GetString("unitname")
	p = order.Product{Id: id}
	p.ProName = html.EscapeString(proname)
	p.UnitName = html.EscapeString(unitname)
	p.Price, _ = this.GetInt64("price")
	p.Cost, _ = this.GetInt64("cost")
	_, err = order.ProductEdit(p)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "error"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
		return
	}

}

type NextPageUrl struct {
	Pages   int64  `json:"pages"`
	Keyword string `json:"keyword"`
	Cateid  int64  `json:"cateid"`
	Storeid int64  `json:"storeid"`
}

//商品列表
type ProductList struct {
	controllers.StaffController
}

func (this *ProductList) Get() {
	store_id, _ := this.GetInt64("storeid")
	pages, _ := this.GetInt64("pages")
	keyword := html.EscapeString(this.GetString("keyword"))
	cate_id, _ := this.GetInt64("cateid")
	if pages == 0 {
		pages = 1
	}

	if store_id == 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择店铺"}
		this.ServeJSON()
		return
	}
	//店铺所有权鉴定

	storeInfo, err := order.GetStoreById(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择正确的店铺"}
		this.ServeJSON()
		return
	}
	if storeInfo.Uid != this.Uid {
		if !this.IsStorePower(store_id) {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
			this.ServeJSON()
			return
		}
	}
	_, list, err := order.ProductListPages(store_id, cate_id, keyword, pages)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "error"}
		this.ServeJSON()
		return
	} else {
		url := &NextPageUrl{}
		url.Pages = pages + 1
		url.Keyword = keyword
		url.Cateid = cate_id
		url.Storeid = store_id

		res := map[string]interface{}{"data": list, "next": url}
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": res}
		this.ServeJSON()
		return
	}
}
