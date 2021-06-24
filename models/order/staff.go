package order

import (
	"encoding/json"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Staff struct {
	Id       int64
	Uid      int64
	MerId    int64
	StoreId  int64
	Creatime int64
}

func (this *Staff) TableName() string {
	return models.TableName("Staff")
}
func init() {
	orm.RegisterModel(new(Staff))
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
