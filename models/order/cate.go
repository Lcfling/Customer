package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
)

type GoodsCate struct {
	Id       int64
	ParentId int64
	Name     string
	Sort     int
	Child    []GoodsCate `orm:"-"`
}

func (this *GoodsCate) TableName() string {
	return models.TableName("goods_cate")
}
func init() {
	orm.RegisterModel(new(GoodsCate))
}
func GetTree() ([]GoodsCate, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_product"))
	var g []GoodsCate
	_, err := qs.OrderBy("-id").All(&g)
	return g, err
}
