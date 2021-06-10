package order

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Store struct {
	Id       int64
	Uid      int64
	Location string
	Lng      string
	Lat      string
	Name     string
	Address  string
	Phone    string
	Rate     int64
	Closed   int64
	Creatime int64
}

func (this *Store) TableName() string {
	return models.TableName("store")
}
func init() {
	orm.RegisterModel(new(Store))
}

func GetStoreById(id int64) (Store, error) {
	o := orm.NewOrm()
	var s Store
	s = Store{Id: id}
	err := o.Read(&s)
	return s, err

}
func Createstore(store Store) (int64, error) {
	o := orm.NewOrm()
	s := new(Store)
	s.Uid = store.Uid
	if s.Location != "" {

		locs := strings.Split(store.Location, ",")
		if len(locs) == 2 {
			s.Location = store.Location
			s.Lng = locs[0]
			s.Lat = locs[1]
		}

	}
	s.Name = store.Name
	s.Address = store.Address
	s.Rate = store.Rate
	s.Creatime = time.Now().Unix()
	id, err := o.Insert(s)
	return id, err
}
func UpdateStore(store Store) error {

	o := orm.NewOrm()
	s := Store{Id: store.Id}
	err := o.Read(&s)
	if err != nil {
		return err
	}
	if store.Name != "" {
		s.Name = store.Name
	}
	if store.Location != "" {

		locs := strings.Split(store.Location, ",")
		if len(locs) == 2 {
			s.Location = store.Location
			s.Lng = locs[0]
			s.Lat = locs[1]
		}

	}
	if store.Phone != "" {
		s.Phone = store.Phone
	}
	if store.Address != "" {
		s.Address = store.Address
	}
	_, err = o.Update(&s, "name", "location", "phone", "address", "lng", "lat")
	return err
}
func GetStoreList(uid int64) (int64, []Store, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("store"))
	cond := orm.NewCondition()
	if uid != 0 {
		cond = cond.And("uid", uid)
	}

	qs = qs.SetCond(cond)
	var s []Store
	num, err := qs.OrderBy("-id").All(&s)
	return num, s, err
}
func UpdateClosed(store_id int64) (int64, error) {
	o := orm.NewOrm()

	Store := Store{Id: store_id}

	err := o.Read(&Store, "id")
	if nil != err {
		return 0, err
	}

	if Store.Closed == 0 {
		Store.Closed = 1
	} else {
		Store.Closed = 0
	}

	_, err = o.Update(&Store, "closed")
	if err != nil {
		return 0, err
	}
	return Store.Closed, err
}
