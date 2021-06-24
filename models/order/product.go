package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type Product struct {
	Id       int64
	StoreId  int64
	Image    string
	ProName  string
	ProInfo  string
	Keyword  string
	Bigcode  string //大码 外包装的条形码
	Exchange int64  //大码兑换小码的数量
	BarCode  string
	CateId   int64
	Cost     int64
	Price    int64
	UnitName string
	Sort     int64
	Sales    int64
	Stock    int64
	AddTime  int64
	IsDel    int
	IsBig    int `orm:"-"`
}

func (this *Product) TableName() string {
	return models.TableName("store_product")
}
func init() {
	orm.RegisterModel(new(Product))
}

func GetProductByCode(code int64) Product {
	barcode := strconv.FormatInt(code, 10)
	var pro Product
	o := orm.NewOrm()
	pro = Product{BarCode: barcode}
	o.Read(&pro, "bar_code")
	return pro
}

func GetProductById(Id int64) (Product, error) {
	var pro Product
	o := orm.NewOrm()
	pro = Product{Id: Id}
	err := o.Read(&pro, "id")
	return pro, err
}
func GetProductByStoreCode(code int64, store_id int64) Product {
	barcode := strconv.FormatInt(code, 10)
	var pro Product
	o := orm.NewOrm()
	pro = Product{BarCode: barcode, StoreId: store_id}
	o.Read(&pro, "bar_code", "store_id")
	return pro
}
func GetProductByBigCode(code int64, store_id int64) (Product, error) {
	barcode := strconv.FormatInt(code, 10)
	var pro Product
	o := orm.NewOrm()
	pro = Product{Bigcode: barcode, StoreId: store_id}
	err := o.Read(&pro, "bigcode", "store_id")
	return pro, err
}

func ProductAdd(proI Product) (int64, error) {
	o := orm.NewOrm()
	pro := new(Product)

	pro.StoreId = proI.StoreId
	pro.Image = proI.Image
	pro.ProName = proI.ProName
	pro.ProInfo = proI.ProInfo
	pro.Keyword = proI.Keyword
	pro.Bigcode = proI.Bigcode
	pro.Exchange = proI.Exchange
	pro.BarCode = proI.BarCode
	pro.CateId = proI.CateId
	pro.Cost = proI.Cost
	pro.Price = proI.Price
	pro.UnitName = proI.UnitName
	pro.Sort = proI.Sort
	pro.Sales = 0
	pro.Stock = 0

	pro.AddTime = time.Now().Unix()
	pro.IsDel = 0
	id, err := o.Insert(pro)
	return id, err
}
func ProductList(store_id int64, keyword string, lastid int64) (int64, []Product, error) {

	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_product"))
	cond := orm.NewCondition()
	if store_id != 0 {
		cond = cond.And("store_id", store_id)
	}
	if keyword != "" {
		cond = cond.And("pro_name__icontains", keyword)
	}
	if lastid != 0 {
		cond = cond.And("id__lt", lastid)
	}

	qs = qs.SetCond(cond)
	var p []Product
	num, err := qs.Limit(15).OrderBy("-id").All(&p)
	return num, p, err
}
func ProductListPages(store_id, cate_id int64, keyword string, pages int64) (int64, []Product, error) {

	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store_product"))
	cond := orm.NewCondition()
	if store_id != 0 {
		cond = cond.And("store_id", store_id)
	}
	if keyword != "" {
		cond = cond.And("pro_name__icontains", keyword)
	}
	if pages != 0 {
		cond = cond.And("id__lt", pages)
	}
	if cate_id != 0 {
		cond = cond.And("cate_id", cate_id)
	}

	qs = qs.SetCond(cond)
	start := (pages - 1) * 15
	var p []Product
	num, err := qs.Limit(15, start).OrderBy("-sales").All(&p)
	return num, p, err
}

func ProductEdit(proI Product) (int64, error) {
	o := orm.NewOrm()
	pro := Product{Id: proI.Id}

	err := o.Read(&pro, "id")
	if proI.ProName != "" {
		pro.ProName = proI.ProName
	}
	if proI.UnitName != "" {
		pro.UnitName = proI.UnitName
	}
	if proI.Cost != 0 {
		pro.Cost = proI.Cost
	}
	if proI.Price != 0 {
		pro.Price = proI.Price
	}
	nums, err := o.Update(&pro, "pro_name", "unit_name", "cost", "price")
	return nums, err
}

/*func GetProList([]int64)([]Product){

}*/
