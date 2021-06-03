package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

//销售详情列表
type Settlement struct {
	Id 				int64
	Day 			string
	Sellcounts 		int64
	Ordercounts 	int64
	Ordersuccess 	int64
	Visitcounts 	int64
	Trademoney		int64
	Ordermoney 		int64
	Uptime 			int64
}

func (this *Settlement) TableName() string {
	return models.TableName("settlement")
}
func init() {
	orm.RegisterModel(new(Settlement))
}

func SaveSettlement(s *Settlement)error{
	o := orm.NewOrm()
	o.Using("default")
	s.Uptime=time.Now().Unix()
	_,err:=o.Insert(s)
	return err
}