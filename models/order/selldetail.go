package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

//销售详情列表
type SellDetail struct {
	Id          int64
	ProductId   int64
	ProductName string
	OrderId     string
	StoreId     int64
	Uid         int64
	Sid         int64
	Nums        float64
	Price       int64
	TotalPrice  int64
	TotalCost   int64
	Status      int
	PayType     int64
	Creatime    int64
}

func (this *SellDetail) TableName() string {
	return models.TableName("sell_detail")
}
func init() {
	orm.RegisterModel(new(SellDetail))
}

func ListByOrder(ordersn string) (int64, []SellDetail, error) {

	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("sell_detail"))
	cond := orm.NewCondition()
	cond = cond.And("order_id", ordersn)
	qs = qs.SetCond(cond)
	var selllist []SellDetail
	num, err := qs.All(&selllist)
	return num, selllist, err

}

func GetSellCounts(b, e int64) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sum(nums) as counts").From("eb_sell_detail").
		Where("creatime>=? and creatime<?")
	sql := qb.String()
	var Sumcounts Sumsells
	err := o.Raw(sql, b, e).QueryRow(&Sumcounts)
	if err != nil {
		return 0, err
	}
	return Sumcounts.Counts, nil
}

//销售 成本
func GetStoreSellCosts(store_id, b, e int64) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	if e == 0 {
		e = time.Now().Unix()
	}
	qb.Select("sum(total_cost) as counts").From("eb_sell_detail").
		Where("store_id=? and creatime>=? and creatime<? and status=1")
	sql := qb.String()
	var sellcost Sumsells
	err := o.Raw(sql, store_id, b, e).QueryRow(&sellcost)
	if err != nil {
		return 0, err
	}
	return sellcost.Counts, nil
}

//销售金额
func GetStoreSellPrice(store_id, b, e int64) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	if e == 0 {
		e = time.Now().Unix()
	}
	qb.Select("sum(total_price) as counts").From("eb_sell_detail").
		Where("store_id=? and creatime>=? and creatime<? and status=1")
	sql := qb.String()
	var sellcost Sumsells
	err := o.Raw(sql, store_id, b, e).QueryRow(&sellcost)
	if err != nil {
		return 0, err
	}
	return sellcost.Counts, nil
}
