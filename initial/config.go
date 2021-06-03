package initial

import (
	"github.com/Lcfling/Customer/models/config"
	"github.com/astaxie/beego"
	"time"
)

func Intconfig(){

	time.Sleep(time.Duration(1)*time.Second)
	c,err:=config.GetConfigList()
	if err!=nil{

	}
	//configMaps:=make(map[string]string)
	for _,v:=range c{
		//configMaps[v.Type]=v.Value
		beego.AppConfig.Set(v.Type,v.Value)
	}
	go Inttask()
}