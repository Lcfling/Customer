package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
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
	Status      int
	PayType     int64
	Creatime    int64
}
type Sumsells struct {
	Counts int64
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
