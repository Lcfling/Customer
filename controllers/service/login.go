package service

import (
	"fmt"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"strconv"
)

type ServiceLogin struct {
	controllers.IndexController
}

func (this *ServiceLogin) Post() {
	account := this.GetString("account")
	passwrod := this.GetString("password")
	if passwrod == "" || account == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "账号密码不能为空"}
		this.ServeJSON()
		return
	}
	user, err := users.GetUserByAccount(account)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "用户不存在"}
		this.ServeJSON()
		return
	}
	if user.Pwd != utils.Md5(passwrod) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "密码错误"}
		this.ServeJSON()
		return
	}
	if user.UserType != 4 && user.UserType != 5 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有登录权限"}
		this.ServeJSON()
		return
	}

	token := utils.Md5(strconv.FormatInt(utils.SnowFlakeId(), 10))
	err = users.UpdateToken(user.Id, token)
	mark := "客服登录"
	go logs.AddUserlog(user.Id, mark, 1, user.UserType)
	user.Token = token
	if err != nil {
		fmt.Println("登录失败：", err.Error())
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败，系统错误！"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功", "data": user}
		this.ServeJSON()
		return
	}
}

type Mqstatus struct {
	controllers.IndexController
}

func (this *Mqstatus) Get() {
	//utils.DEL(models.GetRedis(),"HIGHMQ")
	serlist, _ := utils.LRange(models.GetRedis(), "HIGHMQ", 0, -1)
	fmt.Println(serlist)
	for _, v := range serlist {
		fmt.Println(string(v.([]uint8)))
	}
	serlist2, _ := utils.LRange(models.GetRedis(), "STARTMQ", 0, -1)
	fmt.Println(serlist2)
	for _, v := range serlist2 {
		fmt.Println(string(v.([]uint8)))
	}
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功", "HIGHMQ": serlist, "STARTMQ": serlist2}
	this.ServeJSON()
}
