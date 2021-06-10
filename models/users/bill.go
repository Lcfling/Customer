package users

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
)

type UserBill struct {
	Id      int64
	Uid     int64
	StoreId int64
	OrderId string
	Type    int
	Balance int64
	Mark    string
	AddTime int64
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
	return bill, err
}
