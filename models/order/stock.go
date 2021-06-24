package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type Stock struct {
	Id       int64
	BarCode  int64
	Uid      int64
	Pid      int64
	Pname    string
	StoreId  int64
	Counts   int64
	Creatime int64
}

func (this *Stock) TableName() string {
	return models.TableName("stock_detail")
}
func init() {
	orm.RegisterModel(new(Stock))
}
func AddStock(id, nums, uid int64) error {
	pro, err := GetProductById(id)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Update("eb_sotre_product").Set("stock=stock+?").Where("id=?")
	sql := qb.String()
	o.Begin()
	_, err = o.Raw(sql, nums, id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	qb, _ = orm.NewQueryBuilder("mysql")
	qb.InsertInto("eb_stock_detail", "bar_code", "pname", "uid", "pid", "store_id", "counts", "creatime").Values("?,?,?,?,?,?,?")
	sql = qb.String()
	_, err = o.Raw(sql, pro.BarCode, pro.ProName, uid, id, pro.StoreId, nums, time.Now().Unix()).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	return nil
}
