package order

import (
	"errors"
	"fmt"
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type Order struct {
	Id             int64
	OrderId        string
	OutTradeSn     string
	Uid            int64
	TotalNum       int64
	TotalPrice     int64
	PayPrice       int64
	DeductionPrice int64
	CouponId       int64
	CouponPrice    int64
	Paid           int
	PayTime        int64
	PayType        int64
	AddTime        int64
	Status         int
	Mark           string
	IsDel          int
	Remark         string
	MerId          int64
	StoreId        int64
	Eid            int64
}

//用户提交数据列表
type SubList struct {
	Productid int64
	Nums      int64
}

func (this *Order) TableName() string {
	return models.TableName("store_order")
}
func init() {
	orm.RegisterModel(new(Order))
}

func GetOrdersByUid(uid int64, lastid int64) ([]Order, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_order"))
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.And("uid", uid))
	cond = cond.AndCond(cond.And("status", 1))
	if lastid != 0 {
		cond = cond.AndCond(cond.And("id__lt", lastid))
	}
	qs = qs.SetCond(cond)
	var orderList []Order
	_, err := qs.OrderBy("-id").Limit(10).All(&orderList)
	return orderList, err
}

//按storeid查询店铺销售列表
func GetOrdersBySId(storeid int64, lastid int64) ([]Order, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_order"))
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.And("store_id", storeid))
	if lastid != 0 {
		cond = cond.AndCond(cond.And("id__lt", lastid))
	}
	qs = qs.SetCond(cond)
	var orderList []Order
	_, err := qs.OrderBy("-id").Limit(10).All(&orderList)
	return orderList, err
}

func GetOrderByOrderId(order_id string) (Order, error) {
	o := orm.NewOrm()
	order := Order{OrderId: order_id}
	err := o.Read(&order, "order_id")
	return order, err
}

func CreatOrder(list []SubList, store_id int64, uid int64, order_id string, paytype, eid int64) error {

	//var maplist map[int64]SubList
	maplist := make(map[int64]SubList)
	ids := []int64{}
	for _, v := range list {
		maplist[v.Productid] = v
		ids = append(ids, v.Productid)
	}
	StoreInfo, err := GetStoreById(store_id)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_product"))
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.And("id__in", ids))
	qs = qs.SetCond(cond)
	var Pro []Product
	qs = qs.OrderBy("-id")
	_, err = qs.All(&Pro) //num  总计多少种商品
	if err != nil {
		return err
	}

	var selldetaillist []SellDetail

	var sums int64 = 0       //所有物品的总数量
	var totalprice int64 = 0 //所以物品的总价格
	for _, v := range Pro {
		//selldetaillist
		var sellArr SellDetail
		sellArr.ProductId = v.Id
		sellArr.ProductName = v.ProName
		sellArr.OrderId = order_id
		sellArr.StoreId = store_id
		sellArr.Uid = uid
		sellArr.Nums = maplist[v.Id].Nums
		sellArr.Price = v.Price
		sellArr.TotalPrice = maplist[v.Id].Nums * v.Price
		sellArr.Status = 0
		sellArr.PayType = paytype
		sellArr.Creatime = time.Now().Unix()
		sums += sellArr.Nums
		totalprice += sellArr.TotalPrice
		selldetaillist = append(selldetaillist, sellArr)
	}

	//开启事务  开始入库
	err = o.Begin()
	if err != nil {
		return err
	}
	successNums, err := o.InsertMulti(len(selldetaillist), selldetaillist)
	if len(selldetaillist) != int(successNums) {
		o.Rollback()
		return err
	}

	// 入库订单号
	order := new(Order)
	order.Uid = uid
	order.OrderId = order_id
	order.Mark = ""
	order.TotalNum = sums
	order.TotalPrice = totalprice
	order.PayPrice = 0
	order.DeductionPrice = 0
	order.CouponId = 0
	order.CouponPrice = 0
	order.Paid = 0
	order.PayTime = 0
	order.PayType = paytype
	order.PayPrice = 0
	order.AddTime = time.Now().Unix()
	order.Status = 0
	order.MerId = StoreInfo.Uid
	order.StoreId = store_id
	order.Eid = eid
	_, err = o.Insert(order)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

/*func UpdateOutTradeSn(order_id,out_trade_sn string)error{
	o := orm.NewOrm()
	order:=Order{OrderId:order_id}
	err:=o.Read(&order,"order_id")
	if err!=nil{
		return err
	}
	order.OutTradeSn=out_trade_sn
	_,err=o.Update(&order,"out_trade_sn")
	return err
}*/
//处理订单成功支付
func OrderPaid(order_id string, paidPrice int64, out_trade_sn string) error {

	qb, _ := orm.NewQueryBuilder("mysql")

	qb.Select("*").
		From("eb_store_order").
		Where("order_id=? ").
		And("status=0").
		ForUpdate()
	sql := qb.String()
	o := orm.NewOrm()
	o.Using("default")
	o.Begin()

	var order Order

	err := o.Raw(sql, order_id).QueryRow(&order)

	if err != nil {
		o.Rollback()
		return err
	}
	if order.Status != 0 {
		o.Rollback()
		return errors.New("订单已处理")
	}
	order.PayPrice = paidPrice
	order.PayTime = time.Now().Unix()
	order.Status = 1
	order.OutTradeSn = out_trade_sn
	_, err = o.Update(&order, "pay_price", "pay_time", "status", "out_trade_sn")
	if err != nil {
		o.Rollback()
		return err
	}
	storeInfo, err := GetStoreById(order.StoreId)
	if err != nil {
		o.Rollback()
		return err
	}
	rate := storeInfo.Rate

	increament := order.TotalPrice * (10000 - rate) / 10000

	qb, _ = orm.NewQueryBuilder("mysql")
	qb.Update("eb_user").Set("balance=balance+?").Where("uid=?")
	sql = qb.String()
	_, err = o.Raw(sql, increament, order.MerId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	qb, _ = orm.NewQueryBuilder("mysql")
	qb.InsertInto("eb_user_bill", "uid", "store_id", "order_id", "type", "balance", "mark", "add_time").Values("?,?,?,?,?,?,?")
	sql = qb.String()
	mark := fmt.Sprintf("商户收益 扣除%d%% 手续费", rate/100)
	_, err = o.Raw(sql, order.MerId, order.StoreId, order.OrderId, 1, increament, mark, time.Now().Unix()).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	qb, _ = orm.NewQueryBuilder("mysql")
	qb.Update("eb_sell_detail").Set("status=1").Where("order_id=?")
	sql = qb.String()
	_, err = o.Raw(sql, order_id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

//线下支付
func OrderPaidUnderLine(order_id string, paidPrice int64, out_trade_sn string) error {

	qb, _ := orm.NewQueryBuilder("mysql")

	qb.Select("*").
		From("eb_store_order").
		Where("order_id=? ").
		And("status=0").
		ForUpdate()
	sql := qb.String()
	o := orm.NewOrm()
	o.Using("default")
	o.Begin()

	var order Order

	err := o.Raw(sql, order_id).QueryRow(&order)

	if err != nil {
		o.Rollback()
		return err
	}
	if order.Status != 0 {
		o.Rollback()
		return errors.New("订单已处理")
	}
	order.PayPrice = paidPrice
	order.PayTime = time.Now().Unix()
	order.Status = 1
	order.OutTradeSn = out_trade_sn
	_, err = o.Update(&order, "pay_price", "pay_time", "status", "out_trade_sn")
	if err != nil {
		o.Rollback()
		return err
	}

	qb, _ = orm.NewQueryBuilder("mysql")
	qb.Update("eb_sell_detail").Set("status=1").Where("order_id=?")
	sql = qb.String()
	_, err = o.Raw(sql, order_id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func OrderCounts(b, e int64, status int) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_order"))
	cond := orm.NewCondition()
	cond = cond.And("add_time__gte", b)
	cond = cond.And("add_time__lt", e)
	if status != -2 {
		cond = cond.And("status", status)
	}

	qs = qs.SetCond(cond)
	num, err := qs.Count()
	return num, err
}

type Moneytotal struct {
	TotalPrice int64
	PayPrice   int64
}

func GetSellMoneySum(b, e int64, status int, timeclounm string) (Moneytotal, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	statusStr := ""
	if status != -2 {
		statusStr = " and status=" + strconv.Itoa(status)
	}
	qb.Select("sum(total_price) as total_price,sum(pay_price) as pay_price").From("eb_store_order").
		Where(timeclounm + ">=? and " + timeclounm + "<?" + statusStr)
	sql := qb.String()
	var moneytotal Moneytotal
	err := o.Raw(sql, b, e).QueryRow(&moneytotal)
	if err != nil {
		return Moneytotal{}, err
	}
	return moneytotal, nil
}
