package config

import (
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
)
//config_list
type Config struct {
	Id int64
	Title string
	Type 	string
	Style string
	IsSys int64
	Groups 	 string
	Value 	string
	Extra   string
	Note 	string
	Listsort int64
	CreateTime string
	UpdateTime string
}

func (this *Config) TableName() string {
	return models.TableName("config")
}
func init() {
	orm.RegisterModel(new(Config))
}

func GetConfigList() ([]Config,error){
	o:=orm.NewOrm()
	qs := o.QueryTable(models.TableName("config"))

	var c []Config
	_,err:=qs.All(&c)
	return c,err
}