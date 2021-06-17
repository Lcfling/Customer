package store

import (
	"fmt"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"html"
	"strconv"
	"strings"
)

//需要权限
//
type GoodsController struct {
	controllers.MobileController
}

func (this *GoodsController) Get() {
	barcode, _ := this.GetInt64("barcode")
	//store_id,_:=this.GetInt64("storeid")

	pro := order.GetProductByCode(barcode)

	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
	this.ServeJSON()
}
func (this *GoodsController) Post() {

}

//销量查询
type SellCounts struct {
	controllers.MobileController
}

func (this *SellCounts) Get() {
	storeid, _ := this.GetInt64("storeid")
	lastid, _ := this.GetInt64("lastid")
	orderlist, err := order.GetOrdersBySId(storeid, lastid)
	var sellarray [][]order.SellDetail
	for _, v := range orderlist {
		var selllist []order.SellDetail
		_, selllist, _ = order.ListByOrder(v.OrderId)
		sellarray = append(sellarray, selllist)
	}
	data := map[string]interface{}{"orderlist": orderlist, "sellarray": sellarray}
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有更多的信息"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": data}
		this.ServeJSON()
		return
	}
}

//添加商品
type ProductAddController struct {
	controllers.MobileController
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
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
		this.ServeJSON()
		return
	}

	var pro order.Product
	pro.StoreId = store_id
	pro.Image = images
	pro.ProName = proname
	pro.ProInfo = proinfo
	pro.Keyword = keyword
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
	controllers.MobileController
}

func (this *ProductEdit) Get() {
	id, _ := this.GetInt64("id")
	//store_id,_:=this.GetInt64("storeid")
	if !(id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商品id错误"}
		this.ServeJSON()
		return
	}

	pro := order.GetProductById(id)

	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": pro}
	this.ServeJSON()
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

	var p order.Product

	proname := this.GetString("proname")
	unitname := this.GetString("unitname")
	p = order.Product{Id: id}
	p.ProName = html.EscapeString(proname)
	p.UnitName = html.EscapeString(unitname)
	p.Price, _ = this.GetInt64("price")
	p.Cost, _ = this.GetInt64("cost")
	_, err := order.ProductEdit(p)
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

//商品列表
type ProductList struct {
	controllers.MobileController
}

func (this *ProductList) Get() {
	store_id, _ := this.GetInt64("storeid")
	lastid, _ := this.GetInt64("lastid")
	keyword := this.GetString("keyword")
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
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
		this.ServeJSON()
		return
	}
	_, list, err := order.ProductList(store_id, keyword, lastid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "error"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}
}

// 创建店铺
type CreateStore struct {
	controllers.MobileController
}

func (this *CreateStore) Post() {
	Name := this.GetString("name")
	if Name == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺名称不能为空"}
		this.ServeJSON()
		return
	}
	Location := this.GetString("location")
	if Location == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "坐标不能为空"}
		this.ServeJSON()
		return
	}
	Addr := this.GetString("addr")
	if Addr == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "地址不能为空"}
		this.ServeJSON()
		return
	}
	var s order.Store
	s.Name = Name
	s.Address = Addr
	s.Location = Location
	//todo 配置参数中获取费率
	s.Rate = 600
	id, err := order.Createstore(s)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "创建店铺失败，请联系管理员"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": id}
		this.ServeJSON()
		return
	}
}

//订单列表
type OrderList struct {
	controllers.MobileController
}

func (this *OrderList) Get() {
	lastid, _ := this.GetInt64("lastid")
	store_id, _ := this.GetInt64("storeid")
	// 鉴权 store_id
	storeInfo, err := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}
	orderlist, err := order.GetOrdersBySId(store_id, lastid)
	var sellarray [][]order.SellDetail
	for _, v := range orderlist {
		var selllist []order.SellDetail
		_, selllist, _ = order.ListByOrder(v.OrderId)
		sellarray = append(sellarray, selllist)
	}
	data := map[string]interface{}{"orderlist": orderlist, "sellarray": sellarray}
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有更多的信息"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": data}
		this.ServeJSON()
		return
	}
}

type OrderDetail struct {
	controllers.BaseController
}

func (this *OrderDetail) Get() {
	order_id := this.GetString("orderid")
	orderInfo, err := order.GetOrderByOrderId(order_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "订单不存在"}
		this.ServeJSON()
		return
	}
	_, selllist, err := order.ListByOrder(order_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "订单详情为空"}
		this.ServeJSON()
		return
	}
	//todo 拿到视频播放信息
	data := map[string]interface{}{"orderinfo": orderInfo, "selllist": selllist}
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": data}
	this.ServeJSON()
}

type WithDraw struct {
	controllers.MobileController
}

func (this *WithDraw) Post() {
	money, _ := this.GetInt64("money")
	fid, _ := this.GetInt64("fid")
	//types, _ := this.GetInt("type")

	if !(fid > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择正确的提现方式"}
		this.ServeJSON()
		return
	}
	finance, err := users.GetFinanceById(fid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有找到相关账号"}
		this.ServeJSON()
		return
	}
	types := finance.Type
	order_id := utils.GetOrderSN()
	if money < 100 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "最低提现额度为1元"}
		this.ServeJSON()
		return
	}
	err = order.SubWithdraw(order_id, money, this.Uid, fid, int(types))
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "提现成功 请等待审核！"}
		this.ServeJSON()
		return
	}
}

type StoreList struct {
	controllers.MobileController
}

func (this *StoreList) Get() {
	_, list, err := order.GetStoreList(this.Uid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}
}

type StoreInfo struct {
	controllers.MobileController
}

func (this *StoreInfo) Get() {
	store_id, _ := this.GetInt64("storeid")
	storeInfo, err := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}

	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": storeInfo}
		this.ServeJSON()
	}
}
func (this *StoreInfo) Post() {
	store_id, err := this.GetInt64("storeid")

	Name := html.EscapeString(this.GetString("name"))

	Location := html.EscapeString(this.GetString("location"))
	phone := html.EscapeString(this.GetString("phone"))

	Addr := html.EscapeString(this.GetString("addr"))
	if Location == "" && Name == "" && Addr == "" && phone == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "无变更数据"}
		this.ServeJSON()
		return
	}
	var s order.Store
	s = order.Store{Id: store_id}
	if Name != "" {
		s.Name = Name
	}
	if Location != "" {
		s.Location = Location
	}
	if Addr != "" {
		s.Address = Addr
	}
	if phone != "" {
		s.Phone = phone
	}

	storeInfo, err := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}
	err = order.UpdateStore(s)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "更新店铺信息失败，请联系管理员"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
		return
	}
}

type StoreStatus struct {
	controllers.MobileController
}

func (this *StoreStatus) Post() {
	store_id, _ := this.GetInt64("storeid")

	if store_id == 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	}
	storeInfo, err := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}

	status, err := order.UpdateClosed(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺不存在,404"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": status}
		this.ServeJSON()
		return
	}
}

type EnterLogs struct {
	controllers.MobileController
}

func (this *EnterLogs) Get() {
	store_id, _ := this.GetInt64("storeid")
	lastid, _ := this.GetInt64("lastid")

	//order_id:=this.GetString("orderid")
	if store_id == 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	}
	//todo 鉴权
	storeInfo, err := order.GetStoreById(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}

	list, err := logs.GetLogsByStore(store_id, lastid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "filad:" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}
}

type EnterDetail struct {
	controllers.MobileController
}

func (this *EnterDetail) Get() {
	id, _ := this.GetInt64("id")

	if !(id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的出入记录"}
		this.ServeJSON()
		return
	}

	detail, err := logs.GetEnterLogDetail(id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "filad:" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": detail}
		this.ServeJSON()
		return
	}
}

type GetVideos struct {
	controllers.MobileController
}

func (this *GetVideos) Get() {
	store_id, _ := this.GetInt64("storeid")

	//order_id:=this.GetString("orderid")
	if store_id == 0 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "店铺不存在"}
		this.ServeJSON()
		return
	}
	//todo 鉴权
	storeInfo, err := order.GetStoreById(store_id)
	if err != nil || storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}
	list, err := device.GetVideo(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "filad:" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}
}
