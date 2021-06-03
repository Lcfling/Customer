package logs

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type EnterLog struct {
	Id int64
	Uid int64
	EnterTime int64
	LeaveTime int64
	StoreId int64
	DoorId int64
	ServiceId int64 //客服id
	Status int64 //状态
	Creatime int64 //状态
}

type CloseStore struct {
	StoreId int64
	SerivceId int64
}

func (this *EnterLog) TableName() string {
	return models.TableName("enter_log")
}
func init() {
	orm.RegisterModel(new(EnterLog))
}

func AddEnterLog(logs EnterLog) (int64,error){
	o := orm.NewOrm()
	o.Using("default")
	out:=EnterLog{Uid:logs.Uid,StoreId:logs.StoreId,Status:0}
	err := o.Read(&out, "uid","store_id","status")
	if err!=nil{
		res := new(EnterLog)
		res.Uid=logs.Uid
		res.ServiceId=logs.ServiceId
		res.StoreId=logs.StoreId
		res.Status=0
		res.EnterTime=logs.EnterTime
		res.DoorId=logs.DoorId
		res.LeaveTime=0
		res.Creatime=time.Now().Unix()
		return o.Insert(res)
	}else{
		return out.Id,nil
	}
}

func CustomerLeave(uid int64) error{

	o := orm.NewOrm()
	out:=EnterLog{Uid:uid,Status:0}
	err := o.Read(&out, "uid","status")
	if err!=nil{
		return err
	}else{
		out.LeaveTime=time.Now().Unix()
		out.Status=1
		_,err:=o.Update(&out,"leave_time","status")
		return err
	}
}
func CloseStoreById (storeid int64) error{
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("enter_log"))
	qs.Filter("store_id",storeid).Filter("status",0)

	_,err:=qs.Update(orm.Params{
		"leavetime": time.Now().Unix(),
		"status":1,
	})
	return err
}

func UpdateServiceId (serivceid,uid int64)error{
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("enter_log"))
	qs.Filter("uid",uid).Filter("status",0)

	_,err:=qs.Update(orm.Params{
		"service_id": serivceid,

	})
	return err
}
func UpdateOrderid(storeid,uid int64,orderid string)error{
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("enter_log"))
	qs.Filter("store_id",storeid).Filter("status",0).Filter("uid",uid)

	_,err:=qs.Update(orm.Params{
		"order_id": orderid,
	})
	return err
}

func GetEid(uid int64) (EnterLog,error) {
	o := orm.NewOrm()
	out:=EnterLog{Uid:uid,Status:0}
	err := o.Read(&out, "uid","status")
	return out,err
}
func GetCounts(b,e int64)(int64,error){
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("enter_log"))
	cond := orm.NewCondition()
	cond = cond.And("enter_time__gte", b)
	cond = cond.And("enter_time__lt", e)


	qs = qs.SetCond(cond)
	num, err := qs.Count()
	return num,err
}
func GetLogsByStore(store_id ,lastid int64)([]EnterLog,error){
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable(models.TableName("enter_log"))
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.And("store_id", store_id))
	if lastid!=0 {
		cond = cond.AndCond(cond.And("id__lt", lastid))
	}
	qs = qs.SetCond(cond)
	var enterList []EnterLog
	_,err:=qs.OrderBy("-id").Limit(10).All(&enterList)
	return enterList,err
}
func GetEnterLogDetail(id int64)(EnterLog,error){
	o := orm.NewOrm()
	out:=EnterLog{Id:id}
	err := o.Read(&out)
	return out,err
}