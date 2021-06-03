package device

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
)

type Device struct {
	Id int64
	StoreId int64
	Nums 	int8
	Devicesn int64
	Describe string
	Types 	 int
	Token 		string
}

func (this *Device) TableName() string {
	return models.TableName("store_device")
}
func init() {
	orm.RegisterModel(new(Device))
}

func GetDiviveById(door_id int64) (Device,error){
	o:=orm.NewOrm()
	door:=Device{Id:door_id}
	err:=o.Read(&door)
	return door,err
}
func GetDiviveByToken(token string) (Device,error){
	o:=orm.NewOrm()
	door:=Device{Token:token}
	err:=o.Read(&door,"token")
	return door,err
}
func GetLampByStore(store_id int64,types int)(Device,error){
	o:=orm.NewOrm()
	door:=Device{StoreId:store_id,Types:types}
	err:=o.Read(&door,"storeid","types")
	return door,err
}