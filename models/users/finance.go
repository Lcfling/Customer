package users

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type Finance struct {
	Id       int64
	Uid      int64
	Type     int64
	Code     string
	Name     string
	Bankname string
	Creatime int64
}

func (this *Finance) TableName() string {
	return models.TableName("finance_info")
}
func init() {
	orm.RegisterModel(new(Finance))
}

func GetFinanceById(id int64) (Finance, error) {
	o := orm.NewOrm()
	finance := Finance{Id: id}
	err := o.Read(&finance)
	return finance, err
}
func GetFinanceByUid(Uid int64) (Finance, error) {
	o := orm.NewOrm()
	finance := Finance{Id: Uid}
	err := o.Read(&finance)
	return finance, err
}

func AddFinance(proI Finance) (int64, error) {
	o := orm.NewOrm()
	pro := new(Finance)
	pro.Uid = proI.Uid
	pro.Type = proI.Type
	pro.Name = proI.Name
	pro.Code = proI.Code
	pro.Bankname = proI.Bankname
	pro.Creatime = time.Now().Unix()
	id, err := o.Insert(pro)
	return id, err
}
func GetFinanceList(uid int64) ([]Finance, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("finance_info"))
	cond := orm.NewCondition()
	cond = cond.And("uid", uid)
	qs = qs.SetCond(cond)
	var p []Finance
	_, err := qs.OrderBy("-id").All(&p)
	return p, err
}
