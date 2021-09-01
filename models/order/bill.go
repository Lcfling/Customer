package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserBill struct {
	Id        int64
	Uid       int64
	StoreId   string
	OrderId   string
	Type      int
	Balance   int64
	Mark      string
	AddTime   int64
	OrderInfo Order `orm:"-"`
}
type Sumsells struct {
	Counts int64
}

func (this *UserBill) TableName() string {
	return models.TableName("user_bill")
}
func init() {
	orm.RegisterModel(new(UserBill))
}

func BillList(uid, lastid int64) ([]UserBill, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("user_bill"))
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.And("uid", uid))
	if lastid != 0 {
		cond = cond.AndCond(cond.And("id__lt", lastid))
	}
	qs = qs.SetCond(cond)
	var bill []UserBill
	_, err := qs.OrderBy("-id").Limit(10).All(&bill)
	if len(bill) > 0 {
		for k, v := range bill {
			orderinfo, _ := GetOrderByOrderId(v.OrderId)
			bill[k].OrderInfo = orderinfo
		}
	}
	return bill, err
}
func GetIncome(store_id, b, e int64) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	if e == 0 {
		e = time.Now().Unix()
	}
	qb.Select("sum(balance) as counts").From("eb_user_bill").
		Where("store_id=? and add_time>=? and add_time<? and type=1")
	sql := qb.String()
	var sellcost Sumsells
	err := o.Raw(sql, store_id, b, e).QueryRow(&sellcost)
	if err != nil {
		return 0, err
	}
	return sellcost.Counts, nil
}
