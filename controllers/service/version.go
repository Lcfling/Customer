package service

import (
	"github.com/Lcfling/Customer/controllers"
	"time"
)

type Version struct {
	controllers.IndexController
}
func (this *Version)Get(){
	res:=map[string]interface{}{"ver":"1.1.1.0","url":"https://baidu.com","note":"初始版本","date":time.Now().Unix()}
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":res}
	this.ServeJSON()
}
