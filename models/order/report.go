package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type Report struct {
	Id       		int64
	Uid    		int64
	Cid      		int64
	OrderId  		string
	StoreId 		int64
	Mark 		string
	Creatime 		int64
}
func (this *Report) TableName() string {
	return models.TableName("report")
}
func init() {
	orm.RegisterModel(new(Report))
}

func SaveReport(r *Report) (int64,error){
	o := orm.NewOrm()
	r.Creatime=time.Now().Unix()
	id, err := o.Insert(r)
	return id,err
}