package order

import (
	"errors"
	"fmt"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/users"
	"github.com/astaxie/beego/orm"
	"time"
)
type Withdraw struct {
	Id       		int64
	OrderId    		string
	Uid      		int64
	Money	  		int64
	Creatime 		int64
	Uptime	 		int64
	Types 			int
	Status   		int64
	Mark  		 	string
}



func (this *Withdraw) TableName() string {
	return models.TableName("withdraw")
}
func init() {
	orm.RegisterModel(new(Withdraw))
}

func SubWithdraw(order_id string,money int64,uid int64,types int)error{
	o := orm.NewOrm()
	err:=o.Begin()
	if err!=nil{
		return err
	}
	qb,_:=orm.NewQueryBuilder("mysql")
	qb.Select("*").
		From("eb_user").
		Where("uid=?").
		ForUpdate()
	sql:=qb.String()
	var user users.Users
	err=o.Raw(sql, uid).QueryRow(&user)
	if err!=nil{
		o.Rollback()
		return err
	}
	if user.Balance<money{
		o.Rollback()
		return errors.New("余额不足")
	}
	user.Balance=user.Balance-money
	_,err=o.Update(&user,"balance")
	if err!=nil{
		o.Rollback()
		return err
	}
	order:=new(Withdraw)
	order.Uid=uid
	order.OrderId=order_id
	order.Money=money
	order.Creatime=time.Now().Unix()
	order.Types=types
	order.Status=0
	order.Mark=""
	_,err=o.Insert(order)
	if err!=nil{
		o.Rollback()
		return err
	}
	qb,_=orm.NewQueryBuilder("mysql")
	qb.InsertInto("eb_user_bill","uid","store_id","type","balance","mark","add_time").Values("?,?,?,?,?,?")
	sql=qb.String()
	mark:=fmt.Sprintf("商户提现")
	_,err=o.Raw(sql,uid,0,3,-money,mark,time.Now().Unix()).Exec()
	if err!=nil{
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}