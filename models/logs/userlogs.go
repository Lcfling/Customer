package logs

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserLog struct {
	Id int64
	Uid int64
	Mark string
	Usertype int
	Types int64//状态
	Creatime int64
}


func (this *UserLog) TableName() string {
	return models.TableName("user_log")
}
func init() {
	orm.RegisterModel(new(UserLog))
}

func AddUserlog(uid int64,mark string,types int64,usertype int) (int64,error) {

	o:=orm.NewOrm()
	logs:=new(UserLog)
	logs.Uid=uid
	logs.Creatime=time.Now().Unix()
	logs.Types=types
	logs.Usertype=usertype
	logs.Mark=mark
	return o.Insert(logs)
}