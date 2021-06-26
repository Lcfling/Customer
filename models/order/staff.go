package order

import (
	"encoding/json"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type Staff struct {
	Id       int64
	Uid      int64
	MerId    int64
	StoreId  int64
	Creatime int64
}

func (this *Staff) TableName() string {
	return models.TableName("staff")
}

type StaffMer struct {
	Id       int64
	Uid      int64
	MerId    int64
	Creatime int64
}

func (this *StaffMer) TableName() string {
	return models.TableName("staff_mer")
}
func init() {
	orm.RegisterModel(new(Staff), new(StaffMer))
}

func GetStaffListByUid(uid int64) ([]Staff, error) {
	info, err := utils.Get(models.GetRedis(), "Staff:UID:"+strconv.FormatInt(uid, 10))
	if err == nil && info == "" {
		o := orm.NewOrm()
		o.Using("default")
		qs := o.QueryTable(models.TableName("staff"))
		cond := orm.NewCondition()
		cond = cond.And("uid", uid)
		qs = qs.SetCond(cond)
		var p []Staff
		_, err := qs.All(&p)

		ListJson, err := json.Marshal(p)
		if err != nil {
			return p, err
		}

		err = utils.Set(models.GetRedis(), "Staff:UID:"+strconv.FormatInt(uid, 10), string(ListJson))
		return p, err
	} else if info != "" {
		var list []Staff
		err := json.Unmarshal([]byte(info.(string)), &list)
		return list, err
	} else {
		return nil, err
	}
}
func GetStaffListByMerId(mer_id int64) ([]users.Users, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("u.*").From("eb_staff_mer AS t").
		LeftJoin("eb_user AS u").On("u.uid = t.uid").
		Where("t.mer_id=?")
	sql := qb.String()
	var u []users.Users
	_, err := o.Raw(sql, mer_id).QueryRows(&u)
	return u, err
}
func GetStaffBindList(store_id int64) ([]users.Users, error) {
	o := orm.NewOrm()
	o.Using("default")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("u.*").From("eb_staff AS t").
		LeftJoin("eb_user AS u").On("u.uid = t.uid").
		Where("t.store_id=?")
	sql := qb.String()
	var u []users.Users
	_, err := o.Raw(sql, store_id).QueryRows(&u)
	return u, err
}
func AddStaffMer(uid, mer_id int64) error {
	o := orm.NewOrm()
	s := new(StaffMer)
	s.Uid = uid
	s.MerId = mer_id
	s.Creatime = time.Now().Unix()
	_, err := o.Insert(s)
	return err
}
func BindStaff(uid, mer_id, store_id int64) error {
	o := orm.NewOrm()
	s := new(Staff)
	s.Uid = uid
	s.MerId = mer_id
	s.StoreId = store_id
	s.Creatime = time.Now().Unix()
	_, err := o.Insert(s)
	return err
}
func StoreBindStaffList(store_id int64) ([]Staff, error) {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("staff"))
	cond := orm.NewCondition()
	cond = cond.And("store_id", store_id)
	qs = qs.SetCond(cond)
	var p []Staff
	_, err := qs.All(&p)
	return p, err
}
